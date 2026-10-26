import grpc
import warnings
from credit.v1 import credit_pb2 as credit_dot_v1_dot_credit__pb2
GRPC_GENERATED_VERSION = '1.76.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False
try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True
if _version_not_supported:
    raise RuntimeError(f'The grpc package installed is at version {GRPC_VERSION},' + ' but the generated code in credit/v1/credit_pb2_grpc.py depends on' + f' grpcio>={GRPC_GENERATED_VERSION}.' + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}' + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.')

class CompanyServiceStub(object):

    def __init__(self, channel):
        self.CreateCompany = channel.unary_unary('/credit.v1.CompanyService/CreateCompany', request_serializer=credit_dot_v1_dot_credit__pb2.CompanyInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Company.FromString, _registered_method=True)
        self.GetCompany = channel.unary_unary('/credit.v1.CompanyService/GetCompany', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Company.FromString, _registered_method=True)
        self.ListCompanies = channel.unary_unary('/credit.v1.CompanyService/ListCompanies', request_serializer=credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.CompanyList.FromString, _registered_method=True)
        self.UpdateCompany = channel.unary_unary('/credit.v1.CompanyService/UpdateCompany', request_serializer=credit_dot_v1_dot_credit__pb2.CompanyUpdate.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Company.FromString, _registered_method=True)
        self.VerifyCompany = channel.unary_unary('/credit.v1.CompanyService/VerifyCompany', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Company.FromString, _registered_method=True)
        self.ListProducts = channel.unary_unary('/credit.v1.CompanyService/ListProducts', request_serializer=credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.ProductList.FromString, _registered_method=True)
        self.GetProduct = channel.unary_unary('/credit.v1.CompanyService/GetProduct', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Product.FromString, _registered_method=True)
        self.Recommend = channel.unary_unary('/credit.v1.CompanyService/Recommend', request_serializer=credit_dot_v1_dot_credit__pb2.RecommendationInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.RecommendationList.FromString, _registered_method=True)
        self.Calculate = channel.unary_unary('/credit.v1.CompanyService/Calculate', request_serializer=credit_dot_v1_dot_credit__pb2.CalculationInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Calculation.FromString, _registered_method=True)

class CompanyServiceServicer(object):

    def CreateCompany(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetCompany(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListCompanies(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UpdateCompany(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def VerifyCompany(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListProducts(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetProduct(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Recommend(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Calculate(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

def add_CompanyServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {'CreateCompany': grpc.unary_unary_rpc_method_handler(servicer.CreateCompany, request_deserializer=credit_dot_v1_dot_credit__pb2.CompanyInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Company.SerializeToString), 'GetCompany': grpc.unary_unary_rpc_method_handler(servicer.GetCompany, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Company.SerializeToString), 'ListCompanies': grpc.unary_unary_rpc_method_handler(servicer.ListCompanies, request_deserializer=credit_dot_v1_dot_credit__pb2.Empty.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.CompanyList.SerializeToString), 'UpdateCompany': grpc.unary_unary_rpc_method_handler(servicer.UpdateCompany, request_deserializer=credit_dot_v1_dot_credit__pb2.CompanyUpdate.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Company.SerializeToString), 'VerifyCompany': grpc.unary_unary_rpc_method_handler(servicer.VerifyCompany, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Company.SerializeToString), 'ListProducts': grpc.unary_unary_rpc_method_handler(servicer.ListProducts, request_deserializer=credit_dot_v1_dot_credit__pb2.Empty.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.ProductList.SerializeToString), 'GetProduct': grpc.unary_unary_rpc_method_handler(servicer.GetProduct, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Product.SerializeToString), 'Recommend': grpc.unary_unary_rpc_method_handler(servicer.Recommend, request_deserializer=credit_dot_v1_dot_credit__pb2.RecommendationInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.RecommendationList.SerializeToString), 'Calculate': grpc.unary_unary_rpc_method_handler(servicer.Calculate, request_deserializer=credit_dot_v1_dot_credit__pb2.CalculationInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Calculation.SerializeToString)}
    generic_handler = grpc.method_handlers_generic_handler('credit.v1.CompanyService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('credit.v1.CompanyService', rpc_method_handlers)

class CompanyService(object):

    @staticmethod
    def CreateCompany(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/CreateCompany', credit_dot_v1_dot_credit__pb2.CompanyInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Company.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetCompany(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/GetCompany', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Company.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def ListCompanies(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/ListCompanies', credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, credit_dot_v1_dot_credit__pb2.CompanyList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def UpdateCompany(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/UpdateCompany', credit_dot_v1_dot_credit__pb2.CompanyUpdate.SerializeToString, credit_dot_v1_dot_credit__pb2.Company.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def VerifyCompany(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/VerifyCompany', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Company.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def ListProducts(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/ListProducts', credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, credit_dot_v1_dot_credit__pb2.ProductList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetProduct(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/GetProduct', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Product.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Recommend(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/Recommend', credit_dot_v1_dot_credit__pb2.RecommendationInput.SerializeToString, credit_dot_v1_dot_credit__pb2.RecommendationList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Calculate(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.CompanyService/Calculate', credit_dot_v1_dot_credit__pb2.CalculationInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Calculation.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

class ApplicationServiceStub(object):

    def __init__(self, channel):
        self.CreateApplication = channel.unary_unary('/credit.v1.ApplicationService/CreateApplication', request_serializer=credit_dot_v1_dot_credit__pb2.DraftInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.GetApplication = channel.unary_unary('/credit.v1.ApplicationService/GetApplication', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.ListApplications = channel.unary_unary('/credit.v1.ApplicationService/ListApplications', request_serializer=credit_dot_v1_dot_credit__pb2.PageInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.ApplicationPage.FromString, _registered_method=True)
        self.BankApplications = channel.unary_unary('/credit.v1.ApplicationService/BankApplications', request_serializer=credit_dot_v1_dot_credit__pb2.PageInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.ApplicationPage.FromString, _registered_method=True)
        self.PatchApplication = channel.unary_unary('/credit.v1.ApplicationService/PatchApplication', request_serializer=credit_dot_v1_dot_credit__pb2.PatchInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.Validate = channel.unary_unary('/credit.v1.ApplicationService/Validate', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Validation.FromString, _registered_method=True)
        self.Submit = channel.unary_unary('/credit.v1.ApplicationService/Submit', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.Cancel = channel.unary_unary('/credit.v1.ApplicationService/Cancel', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.GetDecision = channel.unary_unary('/credit.v1.ApplicationService/GetDecision', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Decision.FromString, _registered_method=True)
        self.Accept = channel.unary_unary('/credit.v1.ApplicationService/Accept', request_serializer=credit_dot_v1_dot_credit__pb2.AcceptInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.Decline = channel.unary_unary('/credit.v1.ApplicationService/Decline', request_serializer=credit_dot_v1_dot_credit__pb2.DeclineInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.Issue = channel.unary_unary('/credit.v1.ApplicationService/Issue', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Loan.FromString, _registered_method=True)
        self.History = channel.unary_unary('/credit.v1.ApplicationService/History', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.EventList.FromString, _registered_method=True)
        self.GetScoring = channel.unary_unary('/credit.v1.ApplicationService/GetScoring', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Scoring.FromString, _registered_method=True)
        self.Decide = channel.unary_unary('/credit.v1.ApplicationService/Decide', request_serializer=credit_dot_v1_dot_credit__pb2.ManualInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Application.FromString, _registered_method=True)
        self.GetOverview = channel.unary_unary('/credit.v1.ApplicationService/GetOverview', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Overview.FromString, _registered_method=True)
        self.InformationRequests = channel.unary_unary('/credit.v1.ApplicationService/InformationRequests', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.InformationList.FromString, _registered_method=True)
        self.Feed = channel.unary_unary('/credit.v1.ApplicationService/Feed', request_serializer=credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.FeedList.FromString, _registered_method=True)
        self.ListDocuments = channel.unary_unary('/credit.v1.ApplicationService/ListDocuments', request_serializer=credit_dot_v1_dot_credit__pb2.ID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.DocumentList.FromString, _registered_method=True)
        self.UploadDocument = channel.stream_unary('/credit.v1.ApplicationService/UploadDocument', request_serializer=credit_dot_v1_dot_credit__pb2.UploadChunk.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Document.FromString, _registered_method=True)
        self.DeleteDocument = channel.unary_unary('/credit.v1.ApplicationService/DeleteDocument', request_serializer=credit_dot_v1_dot_credit__pb2.DocumentID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.Empty.FromString, _registered_method=True)
        self.DownloadDocument = channel.unary_stream('/credit.v1.ApplicationService/DownloadDocument', request_serializer=credit_dot_v1_dot_credit__pb2.DocumentID.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.DownloadChunk.FromString, _registered_method=True)

class ApplicationServiceServicer(object):

    def CreateApplication(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetApplication(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListApplications(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def BankApplications(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def PatchApplication(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Validate(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Submit(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Cancel(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetDecision(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Accept(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Decline(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Issue(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def History(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetScoring(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Decide(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetOverview(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def InformationRequests(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Feed(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListDocuments(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UploadDocument(self, request_iterator, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DeleteDocument(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DownloadDocument(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

def add_ApplicationServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {'CreateApplication': grpc.unary_unary_rpc_method_handler(servicer.CreateApplication, request_deserializer=credit_dot_v1_dot_credit__pb2.DraftInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'GetApplication': grpc.unary_unary_rpc_method_handler(servicer.GetApplication, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'ListApplications': grpc.unary_unary_rpc_method_handler(servicer.ListApplications, request_deserializer=credit_dot_v1_dot_credit__pb2.PageInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.ApplicationPage.SerializeToString), 'BankApplications': grpc.unary_unary_rpc_method_handler(servicer.BankApplications, request_deserializer=credit_dot_v1_dot_credit__pb2.PageInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.ApplicationPage.SerializeToString), 'PatchApplication': grpc.unary_unary_rpc_method_handler(servicer.PatchApplication, request_deserializer=credit_dot_v1_dot_credit__pb2.PatchInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'Validate': grpc.unary_unary_rpc_method_handler(servicer.Validate, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Validation.SerializeToString), 'Submit': grpc.unary_unary_rpc_method_handler(servicer.Submit, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'Cancel': grpc.unary_unary_rpc_method_handler(servicer.Cancel, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'GetDecision': grpc.unary_unary_rpc_method_handler(servicer.GetDecision, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Decision.SerializeToString), 'Accept': grpc.unary_unary_rpc_method_handler(servicer.Accept, request_deserializer=credit_dot_v1_dot_credit__pb2.AcceptInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'Decline': grpc.unary_unary_rpc_method_handler(servicer.Decline, request_deserializer=credit_dot_v1_dot_credit__pb2.DeclineInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'Issue': grpc.unary_unary_rpc_method_handler(servicer.Issue, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Loan.SerializeToString), 'History': grpc.unary_unary_rpc_method_handler(servicer.History, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.EventList.SerializeToString), 'GetScoring': grpc.unary_unary_rpc_method_handler(servicer.GetScoring, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Scoring.SerializeToString), 'Decide': grpc.unary_unary_rpc_method_handler(servicer.Decide, request_deserializer=credit_dot_v1_dot_credit__pb2.ManualInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Application.SerializeToString), 'GetOverview': grpc.unary_unary_rpc_method_handler(servicer.GetOverview, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Overview.SerializeToString), 'InformationRequests': grpc.unary_unary_rpc_method_handler(servicer.InformationRequests, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.InformationList.SerializeToString), 'Feed': grpc.unary_unary_rpc_method_handler(servicer.Feed, request_deserializer=credit_dot_v1_dot_credit__pb2.Empty.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.FeedList.SerializeToString), 'ListDocuments': grpc.unary_unary_rpc_method_handler(servicer.ListDocuments, request_deserializer=credit_dot_v1_dot_credit__pb2.ID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.DocumentList.SerializeToString), 'UploadDocument': grpc.stream_unary_rpc_method_handler(servicer.UploadDocument, request_deserializer=credit_dot_v1_dot_credit__pb2.UploadChunk.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Document.SerializeToString), 'DeleteDocument': grpc.unary_unary_rpc_method_handler(servicer.DeleteDocument, request_deserializer=credit_dot_v1_dot_credit__pb2.DocumentID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.Empty.SerializeToString), 'DownloadDocument': grpc.unary_stream_rpc_method_handler(servicer.DownloadDocument, request_deserializer=credit_dot_v1_dot_credit__pb2.DocumentID.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.DownloadChunk.SerializeToString)}
    generic_handler = grpc.method_handlers_generic_handler('credit.v1.ApplicationService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('credit.v1.ApplicationService', rpc_method_handlers)

class ApplicationService(object):

    @staticmethod
    def CreateApplication(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/CreateApplication', credit_dot_v1_dot_credit__pb2.DraftInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetApplication(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/GetApplication', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def ListApplications(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/ListApplications', credit_dot_v1_dot_credit__pb2.PageInput.SerializeToString, credit_dot_v1_dot_credit__pb2.ApplicationPage.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def BankApplications(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/BankApplications', credit_dot_v1_dot_credit__pb2.PageInput.SerializeToString, credit_dot_v1_dot_credit__pb2.ApplicationPage.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def PatchApplication(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/PatchApplication', credit_dot_v1_dot_credit__pb2.PatchInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Validate(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Validate', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Validation.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Submit(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Submit', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Cancel(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Cancel', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetDecision(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/GetDecision', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Decision.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Accept(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Accept', credit_dot_v1_dot_credit__pb2.AcceptInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Decline(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Decline', credit_dot_v1_dot_credit__pb2.DeclineInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Issue(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Issue', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Loan.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def History(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/History', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.EventList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetScoring(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/GetScoring', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Scoring.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Decide(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Decide', credit_dot_v1_dot_credit__pb2.ManualInput.SerializeToString, credit_dot_v1_dot_credit__pb2.Application.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def GetOverview(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/GetOverview', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.Overview.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def InformationRequests(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/InformationRequests', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.InformationList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def Feed(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/Feed', credit_dot_v1_dot_credit__pb2.Empty.SerializeToString, credit_dot_v1_dot_credit__pb2.FeedList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def ListDocuments(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/ListDocuments', credit_dot_v1_dot_credit__pb2.ID.SerializeToString, credit_dot_v1_dot_credit__pb2.DocumentList.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def UploadDocument(request_iterator, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.stream_unary(request_iterator, target, '/credit.v1.ApplicationService/UploadDocument', credit_dot_v1_dot_credit__pb2.UploadChunk.SerializeToString, credit_dot_v1_dot_credit__pb2.Document.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def DeleteDocument(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ApplicationService/DeleteDocument', credit_dot_v1_dot_credit__pb2.DocumentID.SerializeToString, credit_dot_v1_dot_credit__pb2.Empty.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

    @staticmethod
    def DownloadDocument(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_stream(request, target, '/credit.v1.ApplicationService/DownloadDocument', credit_dot_v1_dot_credit__pb2.DocumentID.SerializeToString, credit_dot_v1_dot_credit__pb2.DownloadChunk.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)

class ScoringServiceStub(object):

    def __init__(self, channel):
        self.Evaluate = channel.unary_unary('/credit.v1.ScoringService/Evaluate', request_serializer=credit_dot_v1_dot_credit__pb2.ScoringInput.SerializeToString, response_deserializer=credit_dot_v1_dot_credit__pb2.ScoringOutput.FromString, _registered_method=True)

class ScoringServiceServicer(object):

    def Evaluate(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

def add_ScoringServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {'Evaluate': grpc.unary_unary_rpc_method_handler(servicer.Evaluate, request_deserializer=credit_dot_v1_dot_credit__pb2.ScoringInput.FromString, response_serializer=credit_dot_v1_dot_credit__pb2.ScoringOutput.SerializeToString)}
    generic_handler = grpc.method_handlers_generic_handler('credit.v1.ScoringService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('credit.v1.ScoringService', rpc_method_handlers)

class ScoringService(object):

    @staticmethod
    def Evaluate(request, target, options=(), channel_credentials=None, call_credentials=None, insecure=False, compression=None, wait_for_ready=None, timeout=None, metadata=None):
        return grpc.experimental.unary_unary(request, target, '/credit.v1.ScoringService/Evaluate', credit_dot_v1_dot_credit__pb2.ScoringInput.SerializeToString, credit_dot_v1_dot_credit__pb2.ScoringOutput.FromString, options, channel_credentials, insecure, call_credentials, compression, wait_for_ready, timeout, metadata, _registered_method=True)
