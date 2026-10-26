import json
from functools import cache
from pathlib import Path

from risk.features import attributes
from risk.rules_v1 import MAX_SCORE, MIN_SCORE, Factor, Risk, _points_text
from schema import Application, RiskLevel

MODEL_PATH = Path(__file__).parent / "model_v1.json"
NO_REVENUE_POINTS = -200


@cache
def load(path: Path = MODEL_PATH) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def risk_level(score: int, thresholds: dict) -> RiskLevel:
    if score >= thresholds["low_from"]:
        return RiskLevel.LOW
    if score >= thresholds["medium_from"]:
        return RiskLevel.MEDIUM
    return RiskLevel.HIGH


def assess_risk(app: Application, model: dict | None = None) -> Risk:
    model = model or load()
    factors = []
    if app.monthly_revenue_minor <= 0:
        factors.append(Factor("monthly_revenue_minor", "0", NO_REVENUE_POINTS, f"Нет выручки: {_points_text(NO_REVENUE_POINTS)}"))
    for name, (key, text) in attributes(app).items():
        points = model["points"][name].get(key, 0)
        if points:
            factors.append(Factor(name, key, points, f"{text}: {_points_text(points)}"))

    score = max(MIN_SCORE, min(MAX_SCORE, model["base_points"] + sum(f.points for f in factors)))
    factors.sort(key=lambda f: abs(f.points), reverse=True)
    return Risk(score=score, risk_level=risk_level(score, model["thresholds"]), factors=factors, model_version=model["version"])
