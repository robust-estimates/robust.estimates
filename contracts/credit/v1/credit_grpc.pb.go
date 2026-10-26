package creditv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9
const (
	CompanyService_CreateCompany_FullMethodName = "/credit.v1.CompanyService/CreateCompany"
	CompanyService_GetCompany_FullMethodName    = "/credit.v1.CompanyService/GetCompany"
	CompanyService_ListCompanies_FullMethodName = "/credit.v1.CompanyService/ListCompanies"
	CompanyService_UpdateCompany_FullMethodName = "/credit.v1.CompanyService/UpdateCompany"
	CompanyService_VerifyCompany_FullMethodName = "/credit.v1.CompanyService/VerifyCompany"
	CompanyService_ListProducts_FullMethodName  = "/credit.v1.CompanyService/ListProducts"
	CompanyService_GetProduct_FullMethodName    = "/credit.v1.CompanyService/GetProduct"
	CompanyService_Recommend_FullMethodName     = "/credit.v1.CompanyService/Recommend"
	CompanyService_Calculate_FullMethodName     = "/credit.v1.CompanyService/Calculate"
)

type CompanyServiceClient interface {
	CreateCompany(ctx context.Context, in *CompanyInput, opts ...grpc.CallOption) (*Company, error)
	GetCompany(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Company, error)
	ListCompanies(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CompanyList, error)
	UpdateCompany(ctx context.Context, in *CompanyUpdate, opts ...grpc.CallOption) (*Company, error)
	VerifyCompany(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Company, error)
	ListProducts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ProductList, error)
	GetProduct(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Product, error)
	Recommend(ctx context.Context, in *RecommendationInput, opts ...grpc.CallOption) (*RecommendationList, error)
	Calculate(ctx context.Context, in *CalculationInput, opts ...grpc.CallOption) (*Calculation, error)
}
type companyServiceClient struct{ cc grpc.ClientConnInterface }

func NewCompanyServiceClient(cc grpc.ClientConnInterface) CompanyServiceClient {
	return &companyServiceClient{cc}
}
func (c *companyServiceClient) CreateCompany(ctx context.Context, in *CompanyInput, opts ...grpc.CallOption) (*Company, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Company)
	err := c.cc.Invoke(ctx, CompanyService_CreateCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) GetCompany(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Company, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Company)
	err := c.cc.Invoke(ctx, CompanyService_GetCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) ListCompanies(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CompanyList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompanyList)
	err := c.cc.Invoke(ctx, CompanyService_ListCompanies_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) UpdateCompany(ctx context.Context, in *CompanyUpdate, opts ...grpc.CallOption) (*Company, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Company)
	err := c.cc.Invoke(ctx, CompanyService_UpdateCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) VerifyCompany(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Company, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Company)
	err := c.cc.Invoke(ctx, CompanyService_VerifyCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) ListProducts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ProductList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductList)
	err := c.cc.Invoke(ctx, CompanyService_ListProducts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) GetProduct(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Product, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Product)
	err := c.cc.Invoke(ctx, CompanyService_GetProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) Recommend(ctx context.Context, in *RecommendationInput, opts ...grpc.CallOption) (*RecommendationList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecommendationList)
	err := c.cc.Invoke(ctx, CompanyService_Recommend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *companyServiceClient) Calculate(ctx context.Context, in *CalculationInput, opts ...grpc.CallOption) (*Calculation, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Calculation)
	err := c.cc.Invoke(ctx, CompanyService_Calculate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type CompanyServiceServer interface {
	CreateCompany(context.Context, *CompanyInput) (*Company, error)
	GetCompany(context.Context, *ID) (*Company, error)
	ListCompanies(context.Context, *Empty) (*CompanyList, error)
	UpdateCompany(context.Context, *CompanyUpdate) (*Company, error)
	VerifyCompany(context.Context, *ID) (*Company, error)
	ListProducts(context.Context, *Empty) (*ProductList, error)
	GetProduct(context.Context, *ID) (*Product, error)
	Recommend(context.Context, *RecommendationInput) (*RecommendationList, error)
	Calculate(context.Context, *CalculationInput) (*Calculation, error)
	mustEmbedUnimplementedCompanyServiceServer()
}
type UnimplementedCompanyServiceServer struct{}

func (UnimplementedCompanyServiceServer) CreateCompany(context.Context, *CompanyInput) (*Company, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCompany not implemented")
}
func (UnimplementedCompanyServiceServer) GetCompany(context.Context, *ID) (*Company, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompany not implemented")
}
func (UnimplementedCompanyServiceServer) ListCompanies(context.Context, *Empty) (*CompanyList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCompanies not implemented")
}
func (UnimplementedCompanyServiceServer) UpdateCompany(context.Context, *CompanyUpdate) (*Company, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCompany not implemented")
}
func (UnimplementedCompanyServiceServer) VerifyCompany(context.Context, *ID) (*Company, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCompany not implemented")
}
func (UnimplementedCompanyServiceServer) ListProducts(context.Context, *Empty) (*ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}
func (UnimplementedCompanyServiceServer) GetProduct(context.Context, *ID) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedCompanyServiceServer) Recommend(context.Context, *RecommendationInput) (*RecommendationList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recommend not implemented")
}
func (UnimplementedCompanyServiceServer) Calculate(context.Context, *CalculationInput) (*Calculation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Calculate not implemented")
}
func (UnimplementedCompanyServiceServer) mustEmbedUnimplementedCompanyServiceServer() {
}
func (UnimplementedCompanyServiceServer) testEmbeddedByValue() {
}

type UnsafeCompanyServiceServer interface{ mustEmbedUnimplementedCompanyServiceServer() }

func RegisterCompanyServiceServer(s grpc.ServiceRegistrar, srv CompanyServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CompanyService_ServiceDesc, srv)
}
func _CompanyService_CreateCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompanyInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).CreateCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_CreateCompany_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).CreateCompany(ctx, req.(*CompanyInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_GetCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).GetCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_GetCompany_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).GetCompany(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_ListCompanies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).ListCompanies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_ListCompanies_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).ListCompanies(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_UpdateCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompanyUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).UpdateCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_UpdateCompany_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).UpdateCompany(ctx, req.(*CompanyUpdate))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_VerifyCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).VerifyCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_VerifyCompany_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).VerifyCompany(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_ListProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).ListProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_ListProducts_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).ListProducts(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_GetProduct_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).GetProduct(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_Recommend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendationInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).Recommend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_Recommend_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).Recommend(ctx, req.(*RecommendationInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _CompanyService_Calculate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculationInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyServiceServer).Calculate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: CompanyService_Calculate_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyServiceServer).Calculate(ctx, req.(*CalculationInput))
	}
	return interceptor(ctx, in, info, handler)
}

var CompanyService_ServiceDesc = grpc.ServiceDesc{ServiceName: "credit.v1.CompanyService", HandlerType: (*CompanyServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "CreateCompany", Handler: _CompanyService_CreateCompany_Handler}, {MethodName: "GetCompany", Handler: _CompanyService_GetCompany_Handler}, {MethodName: "ListCompanies", Handler: _CompanyService_ListCompanies_Handler}, {MethodName: "UpdateCompany", Handler: _CompanyService_UpdateCompany_Handler}, {MethodName: "VerifyCompany", Handler: _CompanyService_VerifyCompany_Handler}, {MethodName: "ListProducts", Handler: _CompanyService_ListProducts_Handler}, {MethodName: "GetProduct", Handler: _CompanyService_GetProduct_Handler}, {MethodName: "Recommend", Handler: _CompanyService_Recommend_Handler}, {MethodName: "Calculate", Handler: _CompanyService_Calculate_Handler}}, Streams: []grpc.StreamDesc{}, Metadata: "credit/v1/credit.proto"}

const (
	ApplicationService_CreateApplication_FullMethodName   = "/credit.v1.ApplicationService/CreateApplication"
	ApplicationService_GetApplication_FullMethodName      = "/credit.v1.ApplicationService/GetApplication"
	ApplicationService_ListApplications_FullMethodName    = "/credit.v1.ApplicationService/ListApplications"
	ApplicationService_BankApplications_FullMethodName    = "/credit.v1.ApplicationService/BankApplications"
	ApplicationService_PatchApplication_FullMethodName    = "/credit.v1.ApplicationService/PatchApplication"
	ApplicationService_Validate_FullMethodName            = "/credit.v1.ApplicationService/Validate"
	ApplicationService_Submit_FullMethodName              = "/credit.v1.ApplicationService/Submit"
	ApplicationService_Cancel_FullMethodName              = "/credit.v1.ApplicationService/Cancel"
	ApplicationService_GetDecision_FullMethodName         = "/credit.v1.ApplicationService/GetDecision"
	ApplicationService_Accept_FullMethodName              = "/credit.v1.ApplicationService/Accept"
	ApplicationService_Decline_FullMethodName             = "/credit.v1.ApplicationService/Decline"
	ApplicationService_Issue_FullMethodName               = "/credit.v1.ApplicationService/Issue"
	ApplicationService_History_FullMethodName             = "/credit.v1.ApplicationService/History"
	ApplicationService_GetScoring_FullMethodName          = "/credit.v1.ApplicationService/GetScoring"
	ApplicationService_Decide_FullMethodName              = "/credit.v1.ApplicationService/Decide"
	ApplicationService_GetOverview_FullMethodName         = "/credit.v1.ApplicationService/GetOverview"
	ApplicationService_InformationRequests_FullMethodName = "/credit.v1.ApplicationService/InformationRequests"
	ApplicationService_Feed_FullMethodName                = "/credit.v1.ApplicationService/Feed"
	ApplicationService_ListDocuments_FullMethodName       = "/credit.v1.ApplicationService/ListDocuments"
	ApplicationService_UploadDocument_FullMethodName      = "/credit.v1.ApplicationService/UploadDocument"
	ApplicationService_DeleteDocument_FullMethodName      = "/credit.v1.ApplicationService/DeleteDocument"
	ApplicationService_DownloadDocument_FullMethodName    = "/credit.v1.ApplicationService/DownloadDocument"
)

type ApplicationServiceClient interface {
	CreateApplication(ctx context.Context, in *DraftInput, opts ...grpc.CallOption) (*Application, error)
	GetApplication(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error)
	ListApplications(ctx context.Context, in *PageInput, opts ...grpc.CallOption) (*ApplicationPage, error)
	BankApplications(ctx context.Context, in *PageInput, opts ...grpc.CallOption) (*ApplicationPage, error)
	PatchApplication(ctx context.Context, in *PatchInput, opts ...grpc.CallOption) (*Application, error)
	Validate(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Validation, error)
	Submit(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error)
	Cancel(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error)
	GetDecision(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Decision, error)
	Accept(ctx context.Context, in *AcceptInput, opts ...grpc.CallOption) (*Application, error)
	Decline(ctx context.Context, in *DeclineInput, opts ...grpc.CallOption) (*Application, error)
	Issue(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Loan, error)
	History(ctx context.Context, in *ID, opts ...grpc.CallOption) (*EventList, error)
	GetScoring(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Scoring, error)
	Decide(ctx context.Context, in *ManualInput, opts ...grpc.CallOption) (*Application, error)
	GetOverview(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Overview, error)
	InformationRequests(ctx context.Context, in *ID, opts ...grpc.CallOption) (*InformationList, error)
	Feed(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*FeedList, error)
	ListDocuments(ctx context.Context, in *ID, opts ...grpc.CallOption) (*DocumentList, error)
	UploadDocument(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadChunk, Document], error)
	DeleteDocument(ctx context.Context, in *DocumentID, opts ...grpc.CallOption) (*Empty, error)
	DownloadDocument(ctx context.Context, in *DocumentID, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadChunk], error)
}
type applicationServiceClient struct{ cc grpc.ClientConnInterface }

func NewApplicationServiceClient(cc grpc.ClientConnInterface) ApplicationServiceClient {
	return &applicationServiceClient{cc}
}
func (c *applicationServiceClient) CreateApplication(ctx context.Context, in *DraftInput, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_CreateApplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) GetApplication(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_GetApplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) ListApplications(ctx context.Context, in *PageInput, opts ...grpc.CallOption) (*ApplicationPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApplicationPage)
	err := c.cc.Invoke(ctx, ApplicationService_ListApplications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) BankApplications(ctx context.Context, in *PageInput, opts ...grpc.CallOption) (*ApplicationPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApplicationPage)
	err := c.cc.Invoke(ctx, ApplicationService_BankApplications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) PatchApplication(ctx context.Context, in *PatchInput, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_PatchApplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Validate(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Validation, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Validation)
	err := c.cc.Invoke(ctx, ApplicationService_Validate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Submit(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_Submit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Cancel(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_Cancel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) GetDecision(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Decision, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Decision)
	err := c.cc.Invoke(ctx, ApplicationService_GetDecision_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Accept(ctx context.Context, in *AcceptInput, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_Accept_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Decline(ctx context.Context, in *DeclineInput, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_Decline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Issue(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Loan, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Loan)
	err := c.cc.Invoke(ctx, ApplicationService_Issue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) History(ctx context.Context, in *ID, opts ...grpc.CallOption) (*EventList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventList)
	err := c.cc.Invoke(ctx, ApplicationService_History_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) GetScoring(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Scoring, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Scoring)
	err := c.cc.Invoke(ctx, ApplicationService_GetScoring_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Decide(ctx context.Context, in *ManualInput, opts ...grpc.CallOption) (*Application, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Application)
	err := c.cc.Invoke(ctx, ApplicationService_Decide_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) GetOverview(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Overview, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Overview)
	err := c.cc.Invoke(ctx, ApplicationService_GetOverview_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) InformationRequests(ctx context.Context, in *ID, opts ...grpc.CallOption) (*InformationList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InformationList)
	err := c.cc.Invoke(ctx, ApplicationService_InformationRequests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) Feed(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*FeedList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FeedList)
	err := c.cc.Invoke(ctx, ApplicationService_Feed_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) ListDocuments(ctx context.Context, in *ID, opts ...grpc.CallOption) (*DocumentList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DocumentList)
	err := c.cc.Invoke(ctx, ApplicationService_ListDocuments_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) UploadDocument(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadChunk, Document], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ApplicationService_ServiceDesc.Streams[0], ApplicationService_UploadDocument_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadChunk, Document]{ClientStream: stream}
	return x, nil
}

type ApplicationService_UploadDocumentClient = grpc.ClientStreamingClient[UploadChunk, Document]

func (c *applicationServiceClient) DeleteDocument(ctx context.Context, in *DocumentID, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ApplicationService_DeleteDocument_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *applicationServiceClient) DownloadDocument(ctx context.Context, in *DocumentID, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ApplicationService_ServiceDesc.Streams[1], ApplicationService_DownloadDocument_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DocumentID, DownloadChunk]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ApplicationService_DownloadDocumentClient = grpc.ServerStreamingClient[DownloadChunk]
type ApplicationServiceServer interface {
	CreateApplication(context.Context, *DraftInput) (*Application, error)
	GetApplication(context.Context, *ID) (*Application, error)
	ListApplications(context.Context, *PageInput) (*ApplicationPage, error)
	BankApplications(context.Context, *PageInput) (*ApplicationPage, error)
	PatchApplication(context.Context, *PatchInput) (*Application, error)
	Validate(context.Context, *ID) (*Validation, error)
	Submit(context.Context, *ID) (*Application, error)
	Cancel(context.Context, *ID) (*Application, error)
	GetDecision(context.Context, *ID) (*Decision, error)
	Accept(context.Context, *AcceptInput) (*Application, error)
	Decline(context.Context, *DeclineInput) (*Application, error)
	Issue(context.Context, *ID) (*Loan, error)
	History(context.Context, *ID) (*EventList, error)
	GetScoring(context.Context, *ID) (*Scoring, error)
	Decide(context.Context, *ManualInput) (*Application, error)
	GetOverview(context.Context, *ID) (*Overview, error)
	InformationRequests(context.Context, *ID) (*InformationList, error)
	Feed(context.Context, *Empty) (*FeedList, error)
	ListDocuments(context.Context, *ID) (*DocumentList, error)
	UploadDocument(grpc.ClientStreamingServer[UploadChunk, Document]) error
	DeleteDocument(context.Context, *DocumentID) (*Empty, error)
	DownloadDocument(*DocumentID, grpc.ServerStreamingServer[DownloadChunk]) error
	mustEmbedUnimplementedApplicationServiceServer()
}
type UnimplementedApplicationServiceServer struct{}

func (UnimplementedApplicationServiceServer) CreateApplication(context.Context, *DraftInput) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateApplication not implemented")
}
func (UnimplementedApplicationServiceServer) GetApplication(context.Context, *ID) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplication not implemented")
}
func (UnimplementedApplicationServiceServer) ListApplications(context.Context, *PageInput) (*ApplicationPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListApplications not implemented")
}
func (UnimplementedApplicationServiceServer) BankApplications(context.Context, *PageInput) (*ApplicationPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BankApplications not implemented")
}
func (UnimplementedApplicationServiceServer) PatchApplication(context.Context, *PatchInput) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PatchApplication not implemented")
}
func (UnimplementedApplicationServiceServer) Validate(context.Context, *ID) (*Validation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedApplicationServiceServer) Submit(context.Context, *ID) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (UnimplementedApplicationServiceServer) Cancel(context.Context, *ID) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cancel not implemented")
}
func (UnimplementedApplicationServiceServer) GetDecision(context.Context, *ID) (*Decision, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDecision not implemented")
}
func (UnimplementedApplicationServiceServer) Accept(context.Context, *AcceptInput) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Accept not implemented")
}
func (UnimplementedApplicationServiceServer) Decline(context.Context, *DeclineInput) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decline not implemented")
}
func (UnimplementedApplicationServiceServer) Issue(context.Context, *ID) (*Loan, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Issue not implemented")
}
func (UnimplementedApplicationServiceServer) History(context.Context, *ID) (*EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method History not implemented")
}
func (UnimplementedApplicationServiceServer) GetScoring(context.Context, *ID) (*Scoring, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScoring not implemented")
}
func (UnimplementedApplicationServiceServer) Decide(context.Context, *ManualInput) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decide not implemented")
}
func (UnimplementedApplicationServiceServer) GetOverview(context.Context, *ID) (*Overview, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOverview not implemented")
}
func (UnimplementedApplicationServiceServer) InformationRequests(context.Context, *ID) (*InformationList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InformationRequests not implemented")
}
func (UnimplementedApplicationServiceServer) Feed(context.Context, *Empty) (*FeedList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Feed not implemented")
}
func (UnimplementedApplicationServiceServer) ListDocuments(context.Context, *ID) (*DocumentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDocuments not implemented")
}
func (UnimplementedApplicationServiceServer) UploadDocument(grpc.ClientStreamingServer[UploadChunk, Document]) error {
	return status.Errorf(codes.Unimplemented, "method UploadDocument not implemented")
}
func (UnimplementedApplicationServiceServer) DeleteDocument(context.Context, *DocumentID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDocument not implemented")
}
func (UnimplementedApplicationServiceServer) DownloadDocument(*DocumentID, grpc.ServerStreamingServer[DownloadChunk]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadDocument not implemented")
}
func (UnimplementedApplicationServiceServer) mustEmbedUnimplementedApplicationServiceServer() {
}
func (UnimplementedApplicationServiceServer) testEmbeddedByValue() {
}

type UnsafeApplicationServiceServer interface{ mustEmbedUnimplementedApplicationServiceServer() }

func RegisterApplicationServiceServer(s grpc.ServiceRegistrar, srv ApplicationServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ApplicationService_ServiceDesc, srv)
}
func _ApplicationService_CreateApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DraftInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).CreateApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_CreateApplication_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).CreateApplication(ctx, req.(*DraftInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_GetApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).GetApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_GetApplication_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).GetApplication(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_ListApplications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).ListApplications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_ListApplications_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).ListApplications(ctx, req.(*PageInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_BankApplications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).BankApplications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_BankApplications_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).BankApplications(ctx, req.(*PageInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_PatchApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).PatchApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_PatchApplication_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).PatchApplication(ctx, req.(*PatchInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Validate_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Validate(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Submit_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Submit(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Cancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Cancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Cancel_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Cancel(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_GetDecision_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).GetDecision(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_GetDecision_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).GetDecision(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Accept_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Accept(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Accept_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Accept(ctx, req.(*AcceptInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Decline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeclineInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Decline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Decline_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Decline(ctx, req.(*DeclineInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Issue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Issue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Issue_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Issue(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_History_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).History(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_History_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).History(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_GetScoring_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).GetScoring(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_GetScoring_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).GetScoring(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Decide_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ManualInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Decide(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Decide_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Decide(ctx, req.(*ManualInput))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_GetOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).GetOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_GetOverview_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).GetOverview(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_InformationRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).InformationRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_InformationRequests_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).InformationRequests(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_Feed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).Feed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_Feed_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).Feed(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_ListDocuments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).ListDocuments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_ListDocuments_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).ListDocuments(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_UploadDocument_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ApplicationServiceServer).UploadDocument(&grpc.GenericServerStream[UploadChunk, Document]{ServerStream: stream})
}

type ApplicationService_UploadDocumentServer = grpc.ClientStreamingServer[UploadChunk, Document]

func _ApplicationService_DeleteDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DocumentID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).DeleteDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ApplicationService_DeleteDocument_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).DeleteDocument(ctx, req.(*DocumentID))
	}
	return interceptor(ctx, in, info, handler)
}
func _ApplicationService_DownloadDocument_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DocumentID)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ApplicationServiceServer).DownloadDocument(m, &grpc.GenericServerStream[DocumentID, DownloadChunk]{ServerStream: stream})
}

type ApplicationService_DownloadDocumentServer = grpc.ServerStreamingServer[DownloadChunk]

var ApplicationService_ServiceDesc = grpc.ServiceDesc{ServiceName: "credit.v1.ApplicationService", HandlerType: (*ApplicationServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "CreateApplication", Handler: _ApplicationService_CreateApplication_Handler}, {MethodName: "GetApplication", Handler: _ApplicationService_GetApplication_Handler}, {MethodName: "ListApplications", Handler: _ApplicationService_ListApplications_Handler}, {MethodName: "BankApplications", Handler: _ApplicationService_BankApplications_Handler}, {MethodName: "PatchApplication", Handler: _ApplicationService_PatchApplication_Handler}, {MethodName: "Validate", Handler: _ApplicationService_Validate_Handler}, {MethodName: "Submit", Handler: _ApplicationService_Submit_Handler}, {MethodName: "Cancel", Handler: _ApplicationService_Cancel_Handler}, {MethodName: "GetDecision", Handler: _ApplicationService_GetDecision_Handler}, {MethodName: "Accept", Handler: _ApplicationService_Accept_Handler}, {MethodName: "Decline", Handler: _ApplicationService_Decline_Handler}, {MethodName: "Issue", Handler: _ApplicationService_Issue_Handler}, {MethodName: "History", Handler: _ApplicationService_History_Handler}, {MethodName: "GetScoring", Handler: _ApplicationService_GetScoring_Handler}, {MethodName: "Decide", Handler: _ApplicationService_Decide_Handler}, {MethodName: "GetOverview", Handler: _ApplicationService_GetOverview_Handler}, {MethodName: "InformationRequests", Handler: _ApplicationService_InformationRequests_Handler}, {MethodName: "Feed", Handler: _ApplicationService_Feed_Handler}, {MethodName: "ListDocuments", Handler: _ApplicationService_ListDocuments_Handler}, {MethodName: "DeleteDocument", Handler: _ApplicationService_DeleteDocument_Handler}}, Streams: []grpc.StreamDesc{{StreamName: "UploadDocument", Handler: _ApplicationService_UploadDocument_Handler, ClientStreams: true}, {StreamName: "DownloadDocument", Handler: _ApplicationService_DownloadDocument_Handler, ServerStreams: true}}, Metadata: "credit/v1/credit.proto"}

const (
	ScoringService_Evaluate_FullMethodName = "/credit.v1.ScoringService/Evaluate"
)

type ScoringServiceClient interface {
	Evaluate(ctx context.Context, in *ScoringInput, opts ...grpc.CallOption) (*ScoringOutput, error)
}
type scoringServiceClient struct{ cc grpc.ClientConnInterface }

func NewScoringServiceClient(cc grpc.ClientConnInterface) ScoringServiceClient {
	return &scoringServiceClient{cc}
}
func (c *scoringServiceClient) Evaluate(ctx context.Context, in *ScoringInput, opts ...grpc.CallOption) (*ScoringOutput, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ScoringOutput)
	err := c.cc.Invoke(ctx, ScoringService_Evaluate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ScoringServiceServer interface {
	Evaluate(context.Context, *ScoringInput) (*ScoringOutput, error)
	mustEmbedUnimplementedScoringServiceServer()
}
type UnimplementedScoringServiceServer struct{}

func (UnimplementedScoringServiceServer) Evaluate(context.Context, *ScoringInput) (*ScoringOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Evaluate not implemented")
}
func (UnimplementedScoringServiceServer) mustEmbedUnimplementedScoringServiceServer() {
}
func (UnimplementedScoringServiceServer) testEmbeddedByValue() {
}

type UnsafeScoringServiceServer interface{ mustEmbedUnimplementedScoringServiceServer() }

func RegisterScoringServiceServer(s grpc.ServiceRegistrar, srv ScoringServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ScoringService_ServiceDesc, srv)
}
func _ScoringService_Evaluate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScoringInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScoringServiceServer).Evaluate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ScoringService_Evaluate_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScoringServiceServer).Evaluate(ctx, req.(*ScoringInput))
	}
	return interceptor(ctx, in, info, handler)
}

var ScoringService_ServiceDesc = grpc.ServiceDesc{ServiceName: "credit.v1.ScoringService", HandlerType: (*ScoringServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "Evaluate", Handler: _ScoringService_Evaluate_Handler}}, Streams: []grpc.StreamDesc{}, Metadata: "credit/v1/credit.proto"}
