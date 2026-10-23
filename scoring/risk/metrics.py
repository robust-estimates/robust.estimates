"""Метрики качества моделей риска."""

import numpy as np
import pandas as pd


def roc_auc(y_true, risk_score) -> float:
    """ROC AUC; risk_score: чем больше, тем рискованнее. Одинаковые значения учитываются как половина."""
    y = np.asarray(y_true).astype(bool)
    ranks = pd.Series(np.asarray(risk_score, dtype=float)).rank(method="average").to_numpy()
    n_pos, n_neg = y.sum(), (~y).sum()
    if n_pos == 0 or n_neg == 0:
        raise ValueError("Нужны и дефолты, и не-дефолты")
    return float((ranks[y].sum() - n_pos * (n_pos + 1) / 2) / (n_pos * n_neg))


def gini(y_true, risk_score) -> float:
    """Gini = 2·AUC − 1, стандартная метрика различающей способности в кредитном скоринге."""
    return 2 * roc_auc(y_true, risk_score) - 1
