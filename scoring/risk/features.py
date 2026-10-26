from schema import Application, CreditHistory, LegalForm

OKVED_NAMES = {"G": "торговля", "F": "строительство", "C": "производство", "H": "транспорт", "I": "общепит и гостиницы",
               "M": "профессиональные услуги", "J": "IT и связь", "S": "прочие услуги"}

FEATURES = ("business_age", "debt_load", "margin", "company_history", "owner_history", "tax_debt", "bank_client", "industry")


def _business_age(app: Application) -> tuple[str, str]:
    months = app.business_age_months
    key = "<12" if months < 12 else "12-35" if months < 36 else "36-59" if months < 60 else "60+"
    return key, f"Возраст бизнеса {months} мес."


def _debt_load(app: Application) -> tuple[str, str]:
    if app.monthly_revenue_minor <= 0:
        return "no_revenue", "Нет выручки"
    load = app.monthly_debt_payments_minor / app.monthly_revenue_minor
    key = "0" if load == 0 else "0-20" if load <= 0.2 else "20-35" if load <= 0.35 else "35-50" if load <= 0.5 else "50+"
    return key, f"Долговая нагрузка {load:.0%}"


def _margin(app: Application) -> tuple[str, str]:
    if app.monthly_revenue_minor <= 0:
        return "no_revenue", "Маржа не рассчитана"
    margin = (app.monthly_revenue_minor - app.monthly_expenses_minor) / app.monthly_revenue_minor
    key = "<0" if margin < 0 else "0-5" if margin < 0.05 else "5-15" if margin < 0.15 else "15-25" if margin < 0.25 else "25+"
    return key, f"Маржа {margin:.0%}"


def _history(history: CreditHistory | None, checked: bool, whose: str) -> tuple[str, str]:
    if not checked:
        return "not_checked", f"Кредитная история {whose} не запрашивалась"
    if history is None:
        return "none", f"Нет кредитной истории {whose}"
    if history.max_overdue_days_24m >= 90:
        return "90+", f"Просрочка 90+ дней у {whose} за 2 года"
    if history.overdues_30_plus_12m > 0:
        return "30+", f"Просрочки 30+ дней у {whose} за год: {history.overdues_30_plus_12m}"
    return "clean", f"Нет просрочек у {whose}"


def _tax_debt(app: Application) -> tuple[str, str]:
    if app.tax_debt_minor is None:
        return "unknown", "Нет данных ФНС о налоговом долге"
    if app.tax_debt_minor > 0:
        return "yes", f"Налоговый долг {app.tax_debt_minor // 100:,} ₽".replace(",", " ")
    return "no", "Нет налогового долга"


def _bank_client(app: Application) -> tuple[str, str]:
    if app.account_turnover_monthly_minor:
        return "client", "Есть обороты по счёту в банке"
    return "not_client", "Нет оборотов по счёту в банке"


def _industry(app: Application) -> tuple[str, str]:
    section = app.okved_section or "unknown"
    return section, f"Отрасль: {OKVED_NAMES.get(section, 'не указана')}"


def attributes(app: Application) -> dict[str, tuple[str, str]]:
    result = {
        "business_age": _business_age(app),
        "debt_load": _debt_load(app),
        "margin": _margin(app),
        "owner_history": _history(app.owner_credit_history, app.credit_bureau_checked, "собственника"),
        "tax_debt": _tax_debt(app),
        "bank_client": _bank_client(app),
        "industry": _industry(app),
    }
    if app.legal_form == LegalForm.IP:
        result["company_history"] = ("not_applicable", "У ИП кредитная история компании совпадает с историей собственника")
    else:
        result["company_history"] = _history(app.company_credit_history, app.credit_bureau_checked, "компании")
    return result
