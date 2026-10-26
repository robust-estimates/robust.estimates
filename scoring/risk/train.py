import argparse
import json
import math
from pathlib import Path

import numpy as np
import pandas as pd

from data.generator import TARGET_DEFAULT_RATE, generate_companies, to_application
from risk.features import FEATURES, attributes
from risk.metrics import gini
from risk.rules_v1 import assess_risk
from risk.woe import fit_logistic, woe_table

MODEL_VERSION = "model-v1"
MODEL_PATH = Path(__file__).parent / "model_v1.json"

BASE_SCORE = 600
BASE_ODDS = 50
PDO = 50
FACTOR = PDO / math.log(2)
OFFSET = BASE_SCORE - FACTOR * math.log(BASE_ODDS)

LOW_RISK_PD = 0.03
HIGH_RISK_PD = 0.10


def score_for_pd(pd_value: float) -> float:
    return OFFSET + FACTOR * math.log((1 - pd_value) / pd_value)


def pd_for_score(score):
    return 1 / (1 + np.exp((np.asarray(score, dtype=float) - OFFSET) / FACTOR))


def _keys(frame: pd.DataFrame) -> pd.DataFrame:
    apps = [to_application(row) for row in frame.to_dict("records")]
    return pd.DataFrame([{name: key for name, (key, _) in attributes(app).items()} for app in apps])[list(FEATURES)]


def _scores(keys: pd.DataFrame, base_points: int, points: dict) -> np.ndarray:
    total = np.full(len(keys), base_points, dtype=float)
    for name in FEATURES:
        total += keys[name].map(points[name]).fillna(0).to_numpy()
    return np.clip(np.round(total), 0, 1000)


def _levels(scores: np.ndarray, low_from: int, medium_from: int) -> np.ndarray:
    return np.where(scores >= low_from, "LOW", np.where(scores >= medium_from, "MEDIUM", "HIGH"))


def train(n: int = 40_000, seed: int = 42, test_share: float = 0.3) -> dict:
    data = generate_companies(n, seed)
    test = np.random.default_rng(seed).random(n) < test_share
    keys = _keys(data)
    y = data["default_12m"].to_numpy(dtype=float)

    tables = {name: woe_table(keys.loc[~test, name], y[~test]) for name in FEATURES}
    x = np.column_stack([keys[name].map(tables[name]["woe"]).fillna(0).to_numpy() for name in FEATURES])
    weights = fit_logistic(x[~test], y[~test])
    intercept, coefficients = weights[0], dict(zip(FEATURES, weights[1:]))

    base_points = round(OFFSET - FACTOR * intercept)
    points = {name: {key: round(-FACTOR * coefficients[name] * woe) for key, woe in tables[name]["woe"].items()}
              for name in FEATURES}
    low_from, medium_from = math.ceil(score_for_pd(LOW_RISK_PD)), math.ceil(score_for_pd(HIGH_RISK_PD))

    scores = _scores(keys, base_points, points)
    rules = np.array([assess_risk(to_application(row)).score for row in data.loc[test].to_dict("records")])
    levels = _levels(scores[test], low_from, medium_from)
    predicted = pd_for_score(scores[test])
    deciles = pd.qcut(pd.Series(predicted).rank(method="first"), 10, labels=False)
    calibration = pd.DataFrame({"decile": deciles, "predicted": predicted, "observed": y[test]}).groupby("decile").mean()

    return {
        "version": MODEL_VERSION,
        "trained_on": {"generator": "data.generator", "n": n, "seed": seed, "test_share": test_share,
                       "target": "default_12m", "target_default_rate": TARGET_DEFAULT_RATE},
        "scale": {"base_score": BASE_SCORE, "base_odds": BASE_ODDS, "pdo": PDO},
        "base_points": base_points,
        "thresholds": {"low_from": low_from, "medium_from": medium_from, "low_risk_pd": LOW_RISK_PD, "high_risk_pd": HIGH_RISK_PD},
        "points": points,
        "metrics": {
            "gini_test": round(gini(y[test], -scores[test]), 3),
            "gini_train": round(gini(y[~test], -scores[~test]), 3),
            "gini_rules_v1_test": round(gini(y[test], -rules), 3),
            "gini_ceiling_test": round(gini(y[test], data.loc[test, "true_pd"]), 3),
            "coefficients": {name: round(float(value), 3) for name, value in coefficients.items()},
            "information_value": {name: round(float(tables[name]["iv"].sum()), 3) for name in FEATURES},
            "levels_test": {level: {"share": round(float((levels == level).mean()), 3),
                                    "default_rate": round(float(y[test][levels == level].mean()), 3)}
                            for level in ("LOW", "MEDIUM", "HIGH")},
            "calibration_test": [{"predicted_pd": round(float(row.predicted), 3), "observed_rate": round(float(row.observed), 3)}
                                 for row in calibration.itertuples()],
            "bad_rate_by_bin": {name: {key: round(float(rate), 3) for key, rate in tables[name]["bad_rate"].items()}
                                for name in FEATURES},
        },
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="Обучение model-v1: WoE + логистическая регрессия на синтетических данных")
    parser.add_argument("--n", type=int, default=40_000)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--out", type=Path, default=MODEL_PATH)
    args = parser.parse_args()

    model = train(args.n, args.seed)
    args.out.write_text(json.dumps(model, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(model["metrics"], ensure_ascii=False, indent=2))
    print(f"Сохранено в {args.out}")


if __name__ == "__main__":
    main()
