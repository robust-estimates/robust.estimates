import numpy as np
import pandas as pd

SMOOTHING = 0.5


def woe_table(keys, defaults) -> pd.DataFrame:
    frame = pd.DataFrame({"key": np.asarray(keys), "bad": np.asarray(defaults, dtype=float)})
    table = frame.groupby("key")["bad"].agg(count="size", bad="sum")
    table["good"] = table["count"] - table["bad"]
    good_share = (table["good"] + SMOOTHING) / (table["good"].sum() + SMOOTHING * len(table))
    bad_share = (table["bad"] + SMOOTHING) / (table["bad"].sum() + SMOOTHING * len(table))
    table["woe"] = np.log(good_share / bad_share)
    table["bad_rate"] = table["bad"] / table["count"]
    table["iv"] = (good_share - bad_share) * table["woe"]
    return table


def fit_logistic(x: np.ndarray, y: np.ndarray, l2: float = 1e-3, iterations: int = 50) -> np.ndarray:
    x = np.c_[np.ones(len(x)), x]
    weights = np.zeros(x.shape[1])
    penalty = l2 * np.eye(x.shape[1])
    penalty[0, 0] = 0
    for _ in range(iterations):
        p = 1 / (1 + np.exp(-x @ weights))
        hessian = x.T @ (x * (p * (1 - p))[:, None]) + penalty
        step = np.linalg.solve(hessian, x.T @ (y - p) - penalty @ weights)
        weights += step
        if np.abs(step).max() < 1e-10:
            break
    return weights
