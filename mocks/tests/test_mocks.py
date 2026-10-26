import json
import threading
import unittest
import urllib.error
import urllib.request
from http.server import ThreadingHTTPServer

import profiles
from server import Handler


class ProfileTests(unittest.TestCase):
    def test_same_inn_gives_same_answer(self):
        for build in (profiles.registry, profiles.credit_report, profiles.account):
            self.assertEqual(build("7712345678"), build("7712345678"))

    def test_different_inns_give_different_answers(self):
        reports = {json.dumps(profiles.credit_report(f"77{i:08d}"), sort_keys=True) for i in range(50)}
        self.assertGreater(len(reports), 10)

    def test_legal_form_follows_inn_length(self):
        self.assertEqual(profiles.registry("7712345678")["legalForm"], "OOO")
        self.assertEqual(profiles.registry("771234567890")["legalForm"], "IP")

    def test_ip_has_no_separate_company_history(self):
        for i in range(30):
            self.assertIsNone(profiles.credit_report(f"77{i:010d}")["company"])

    def test_invalid_inn_is_rejected(self):
        for inn in ("123", "77123456AB", "1234567890123"):
            with self.assertRaises(profiles.InvalidINN):
                profiles.registry(inn)

    def test_fixed_profiles(self):
        self.assertEqual(profiles.registry("0000090043")["status"], "LIQUIDATED")
        self.assertEqual(profiles.registry("0000090036")["taxDebtMinor"], 12000000)
        self.assertEqual(profiles.credit_report("0000090029")["owner"]["overdueCount12m"], 3)
        self.assertFalse(profiles.credit_report("000000900509")["owner"]["found"])

    def test_e2e_company_is_clean(self):
        report = profiles.credit_report("7700000000")
        self.assertEqual(report["company"]["overdueCount12m"] + report["owner"]["overdueCount12m"], 0)
        self.assertEqual(profiles.registry("7700000000")["taxDebtMinor"], 0)

    def test_generated_values_are_plausible(self):
        for i in range(200):
            inn = f"50{i:08d}"
            fns = profiles.registry(inn)
            self.assertIn(fns["okvedSection"], profiles.OKVED_SECTIONS)
            self.assertGreaterEqual(fns["taxDebtMinor"], 0)
            owner = profiles.credit_report(inn)["owner"]
            self.assertTrue(0 <= owner["overdueCount12m"] <= 10)
            self.assertTrue(0 <= owner["maxOverdueDays24m"] <= 365)
            if owner["overdueCount12m"]:
                self.assertGreaterEqual(owner["maxOverdueDays24m"], 30)


class ServerTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        cls.base = f"http://127.0.0.1:{cls.server.server_port}"
        threading.Thread(target=cls.server.serve_forever, daemon=True).start()

    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown()

    def get(self, path):
        try:
            with urllib.request.urlopen(self.base + path, timeout=2) as reply:
                return reply.status, json.loads(reply.read())
        except urllib.error.HTTPError as error:
            return error.code, json.loads(error.read())

    def test_routes(self):
        self.assertEqual(self.get("/health")[0], 200)
        self.assertEqual(self.get("/fns/v1/companies/7700000000")[1]["status"], "ACTIVE")
        self.assertIn("owner", self.get("/bki/v1/reports/7700000000")[1])
        self.assertTrue(self.get("/bank/v1/accounts/7700000000")[1]["found"])

    def test_errors(self):
        self.assertEqual(self.get("/fns/v1/companies/123")[0], 400)
        self.assertEqual(self.get("/unknown")[0], 404)
        self.assertEqual(self.get("/bki/v1/reports/000000900682")[0], 503)
        self.assertEqual(self.get("/fns/v1/companies/000000900682")[0], 200)


if __name__ == "__main__":
    unittest.main()
