package creditv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_credit_v1_credit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Empty) ProtoMessage() {
}
func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Empty) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{0}
}

type ID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ID) Reset() {
	*x = ID{}
	mi := &file_credit_v1_credit_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ID) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ID) ProtoMessage() {
}
func (x *ID) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ID) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{1}
}
func (x *ID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CompanyInput struct {
	state                    protoimpl.MessageState `protogen:"open.v1"`
	Inn                      string                 `protobuf:"bytes,1,opt,name=inn,proto3" json:"inn,omitempty"`
	Name                     string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	BusinessAgeMonths        int32                  `protobuf:"varint,3,opt,name=business_age_months,json=businessAgeMonths,proto3" json:"business_age_months,omitempty"`
	MonthlyRevenueMinor      int64                  `protobuf:"varint,4,opt,name=monthly_revenue_minor,json=monthlyRevenueMinor,proto3" json:"monthly_revenue_minor,omitempty"`
	MonthlyExpensesMinor     int64                  `protobuf:"varint,5,opt,name=monthly_expenses_minor,json=monthlyExpensesMinor,proto3" json:"monthly_expenses_minor,omitempty"`
	MonthlyDebtPaymentsMinor int64                  `protobuf:"varint,6,opt,name=monthly_debt_payments_minor,json=monthlyDebtPaymentsMinor,proto3" json:"monthly_debt_payments_minor,omitempty"`
	Kpp                      string                 `protobuf:"bytes,7,opt,name=kpp,proto3" json:"kpp,omitempty"`
	Ogrn                     string                 `protobuf:"bytes,8,opt,name=ogrn,proto3" json:"ogrn,omitempty"`
	Industry                 string                 `protobuf:"bytes,9,opt,name=industry,proto3" json:"industry,omitempty"`
	RegisteredAt             string                 `protobuf:"bytes,10,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at,omitempty"`
	RepresentativeName       string                 `protobuf:"bytes,11,opt,name=representative_name,json=representativeName,proto3" json:"representative_name,omitempty"`
	RepresentativePosition   string                 `protobuf:"bytes,12,opt,name=representative_position,json=representativePosition,proto3" json:"representative_position,omitempty"`
	Phone                    string                 `protobuf:"bytes,13,opt,name=phone,proto3" json:"phone,omitempty"`
	Email                    string                 `protobuf:"bytes,14,opt,name=email,proto3" json:"email,omitempty"`
	AnnualRevenueMinor       int64                  `protobuf:"varint,15,opt,name=annual_revenue_minor,json=annualRevenueMinor,proto3" json:"annual_revenue_minor,omitempty"`
	AnnualExpensesMinor      int64                  `protobuf:"varint,16,opt,name=annual_expenses_minor,json=annualExpensesMinor,proto3" json:"annual_expenses_minor,omitempty"`
	CashFlowMinor            int64                  `protobuf:"varint,17,opt,name=cash_flow_minor,json=cashFlowMinor,proto3" json:"cash_flow_minor,omitempty"`
	OutstandingDebtMinor     int64                  `protobuf:"varint,18,opt,name=outstanding_debt_minor,json=outstandingDebtMinor,proto3" json:"outstanding_debt_minor,omitempty"`
	MarketplaceSales         bool                   `protobuf:"varint,20,opt,name=marketplace_sales,json=marketplaceSales,proto3" json:"marketplace_sales,omitempty"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *CompanyInput) Reset() {
	*x = CompanyInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CompanyInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CompanyInput) ProtoMessage() {
}
func (x *CompanyInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CompanyInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{2}
}
func (x *CompanyInput) GetInn() string {
	if x != nil {
		return x.Inn
	}
	return ""
}
func (x *CompanyInput) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *CompanyInput) GetBusinessAgeMonths() int32 {
	if x != nil {
		return x.BusinessAgeMonths
	}
	return 0
}
func (x *CompanyInput) GetMonthlyRevenueMinor() int64 {
	if x != nil {
		return x.MonthlyRevenueMinor
	}
	return 0
}
func (x *CompanyInput) GetMonthlyExpensesMinor() int64 {
	if x != nil {
		return x.MonthlyExpensesMinor
	}
	return 0
}
func (x *CompanyInput) GetMonthlyDebtPaymentsMinor() int64 {
	if x != nil {
		return x.MonthlyDebtPaymentsMinor
	}
	return 0
}
func (x *CompanyInput) GetKpp() string {
	if x != nil {
		return x.Kpp
	}
	return ""
}
func (x *CompanyInput) GetOgrn() string {
	if x != nil {
		return x.Ogrn
	}
	return ""
}
func (x *CompanyInput) GetIndustry() string {
	if x != nil {
		return x.Industry
	}
	return ""
}
func (x *CompanyInput) GetRegisteredAt() string {
	if x != nil {
		return x.RegisteredAt
	}
	return ""
}
func (x *CompanyInput) GetRepresentativeName() string {
	if x != nil {
		return x.RepresentativeName
	}
	return ""
}
func (x *CompanyInput) GetRepresentativePosition() string {
	if x != nil {
		return x.RepresentativePosition
	}
	return ""
}
func (x *CompanyInput) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}
func (x *CompanyInput) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}
func (x *CompanyInput) GetAnnualRevenueMinor() int64 {
	if x != nil {
		return x.AnnualRevenueMinor
	}
	return 0
}
func (x *CompanyInput) GetAnnualExpensesMinor() int64 {
	if x != nil {
		return x.AnnualExpensesMinor
	}
	return 0
}
func (x *CompanyInput) GetCashFlowMinor() int64 {
	if x != nil {
		return x.CashFlowMinor
	}
	return 0
}
func (x *CompanyInput) GetOutstandingDebtMinor() int64 {
	if x != nil {
		return x.OutstandingDebtMinor
	}
	return 0
}
func (x *CompanyInput) GetMarketplaceSales() bool {
	if x != nil {
		return x.MarketplaceSales
	}
	return false
}

type Company struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Id                  string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Data                *CompanyInput          `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	VerificationStatus  string                 `protobuf:"bytes,3,opt,name=verification_status,json=verificationStatus,proto3" json:"verification_status,omitempty"`
	VerificationMessage string                 `protobuf:"bytes,4,opt,name=verification_message,json=verificationMessage,proto3" json:"verification_message,omitempty"`
	OwnerId             string                 `protobuf:"bytes,5,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Version             int64                  `protobuf:"varint,6,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Company) Reset() {
	*x = Company{}
	mi := &file_credit_v1_credit_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Company) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Company) ProtoMessage() {
}
func (x *Company) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Company) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{3}
}
func (x *Company) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Company) GetData() *CompanyInput {
	if x != nil {
		return x.Data
	}
	return nil
}
func (x *Company) GetVerificationStatus() string {
	if x != nil {
		return x.VerificationStatus
	}
	return ""
}
func (x *Company) GetVerificationMessage() string {
	if x != nil {
		return x.VerificationMessage
	}
	return ""
}
func (x *Company) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}
func (x *Company) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type CompanyList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Company             `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompanyList) Reset() {
	*x = CompanyList{}
	mi := &file_credit_v1_credit_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CompanyList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CompanyList) ProtoMessage() {
}
func (x *CompanyList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CompanyList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{4}
}
func (x *CompanyList) GetItems() []*Company {
	if x != nil {
		return x.Items
	}
	return nil
}

type CompanyUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Data          *CompanyInput          `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompanyUpdate) Reset() {
	*x = CompanyUpdate{}
	mi := &file_credit_v1_credit_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CompanyUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CompanyUpdate) ProtoMessage() {
}
func (x *CompanyUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CompanyUpdate) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{5}
}
func (x *CompanyUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *CompanyUpdate) GetData() *CompanyInput {
	if x != nil {
		return x.Data
	}
	return nil
}

type Product struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	Id                    string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                  string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Purposes              []string               `protobuf:"bytes,3,rep,name=purposes,proto3" json:"purposes,omitempty"`
	MinAmountMinor        int64                  `protobuf:"varint,4,opt,name=min_amount_minor,json=minAmountMinor,proto3" json:"min_amount_minor,omitempty"`
	MaxAmountMinor        int64                  `protobuf:"varint,5,opt,name=max_amount_minor,json=maxAmountMinor,proto3" json:"max_amount_minor,omitempty"`
	MinTermMonths         int32                  `protobuf:"varint,6,opt,name=min_term_months,json=minTermMonths,proto3" json:"min_term_months,omitempty"`
	MaxTermMonths         int32                  `protobuf:"varint,7,opt,name=max_term_months,json=maxTermMonths,proto3" json:"max_term_months,omitempty"`
	AnnualRateBps         int32                  `protobuf:"varint,8,opt,name=annual_rate_bps,json=annualRateBps,proto3" json:"annual_rate_bps,omitempty"`
	RequiredDocumentTypes []string               `protobuf:"bytes,9,rep,name=required_document_types,json=requiredDocumentTypes,proto3" json:"required_document_types,omitempty"`
	Description           string                 `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	Collateral            string                 `protobuf:"bytes,11,opt,name=collateral,proto3" json:"collateral,omitempty"`
	PaymentType           string                 `protobuf:"bytes,12,opt,name=payment_type,json=paymentType,proto3" json:"payment_type,omitempty"`
	Tags                  []string               `protobuf:"bytes,13,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *Product) Reset() {
	*x = Product{}
	mi := &file_credit_v1_credit_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Product) ProtoMessage() {
}
func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Product) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{6}
}
func (x *Product) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *Product) GetPurposes() []string {
	if x != nil {
		return x.Purposes
	}
	return nil
}
func (x *Product) GetMinAmountMinor() int64 {
	if x != nil {
		return x.MinAmountMinor
	}
	return 0
}
func (x *Product) GetMaxAmountMinor() int64 {
	if x != nil {
		return x.MaxAmountMinor
	}
	return 0
}
func (x *Product) GetMinTermMonths() int32 {
	if x != nil {
		return x.MinTermMonths
	}
	return 0
}
func (x *Product) GetMaxTermMonths() int32 {
	if x != nil {
		return x.MaxTermMonths
	}
	return 0
}
func (x *Product) GetAnnualRateBps() int32 {
	if x != nil {
		return x.AnnualRateBps
	}
	return 0
}
func (x *Product) GetRequiredDocumentTypes() []string {
	if x != nil {
		return x.RequiredDocumentTypes
	}
	return nil
}
func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}
func (x *Product) GetCollateral() string {
	if x != nil {
		return x.Collateral
	}
	return ""
}
func (x *Product) GetPaymentType() string {
	if x != nil {
		return x.PaymentType
	}
	return ""
}
func (x *Product) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type ProductList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Product             `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductList) Reset() {
	*x = ProductList{}
	mi := &file_credit_v1_credit_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ProductList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ProductList) ProtoMessage() {
}
func (x *ProductList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ProductList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{7}
}
func (x *ProductList) GetItems() []*Product {
	if x != nil {
		return x.Items
	}
	return nil
}

type Need struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CompanyId     string                 `protobuf:"bytes,1,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
	AmountMinor   int64                  `protobuf:"varint,2,opt,name=amount_minor,json=amountMinor,proto3" json:"amount_minor,omitempty"`
	TermMonths    int32                  `protobuf:"varint,3,opt,name=term_months,json=termMonths,proto3" json:"term_months,omitempty"`
	Purpose       string                 `protobuf:"bytes,4,opt,name=purpose,proto3" json:"purpose,omitempty"`
	Description   string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Need) Reset() {
	*x = Need{}
	mi := &file_credit_v1_credit_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Need) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Need) ProtoMessage() {
}
func (x *Need) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Need) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{8}
}
func (x *Need) GetCompanyId() string {
	if x != nil {
		return x.CompanyId
	}
	return ""
}
func (x *Need) GetAmountMinor() int64 {
	if x != nil {
		return x.AmountMinor
	}
	return 0
}
func (x *Need) GetTermMonths() int32 {
	if x != nil {
		return x.TermMonths
	}
	return 0
}
func (x *Need) GetPurpose() string {
	if x != nil {
		return x.Purpose
	}
	return ""
}
func (x *Need) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type RecommendationInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CompanyId     string                 `protobuf:"bytes,1,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
	AmountMinor   *int64                 `protobuf:"varint,2,opt,name=amount_minor,json=amountMinor,proto3,oneof" json:"amount_minor,omitempty"`
	TermMonths    int32                  `protobuf:"varint,3,opt,name=term_months,json=termMonths,proto3" json:"term_months,omitempty"`
	Purpose       string                 `protobuf:"bytes,4,opt,name=purpose,proto3" json:"purpose,omitempty"`
	Collateral    string                 `protobuf:"bytes,5,opt,name=collateral,proto3" json:"collateral,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecommendationInput) Reset() {
	*x = RecommendationInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *RecommendationInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*RecommendationInput) ProtoMessage() {
}
func (x *RecommendationInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*RecommendationInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{9}
}
func (x *RecommendationInput) GetCompanyId() string {
	if x != nil {
		return x.CompanyId
	}
	return ""
}
func (x *RecommendationInput) GetAmountMinor() int64 {
	if x != nil && x.AmountMinor != nil {
		return *x.AmountMinor
	}
	return 0
}
func (x *RecommendationInput) GetTermMonths() int32 {
	if x != nil {
		return x.TermMonths
	}
	return 0
}
func (x *RecommendationInput) GetPurpose() string {
	if x != nil {
		return x.Purpose
	}
	return ""
}
func (x *RecommendationInput) GetCollateral() string {
	if x != nil {
		return x.Collateral
	}
	return ""
}

type Recommendation struct {
	state                        protoimpl.MessageState `protogen:"open.v1"`
	Product                      *Product               `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
	Reason                       string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	EstimatedMonthlyPaymentMinor int64                  `protobuf:"varint,3,opt,name=estimated_monthly_payment_minor,json=estimatedMonthlyPaymentMinor,proto3" json:"estimated_monthly_payment_minor,omitempty"`
	Availability                 string                 `protobuf:"bytes,4,opt,name=availability,proto3" json:"availability,omitempty"`
	Rank                         int32                  `protobuf:"varint,5,opt,name=rank,proto3" json:"rank,omitempty"`
	unknownFields                protoimpl.UnknownFields
	sizeCache                    protoimpl.SizeCache
}

func (x *Recommendation) Reset() {
	*x = Recommendation{}
	mi := &file_credit_v1_credit_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Recommendation) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Recommendation) ProtoMessage() {
}
func (x *Recommendation) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Recommendation) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{10}
}
func (x *Recommendation) GetProduct() *Product {
	if x != nil {
		return x.Product
	}
	return nil
}
func (x *Recommendation) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}
func (x *Recommendation) GetEstimatedMonthlyPaymentMinor() int64 {
	if x != nil {
		return x.EstimatedMonthlyPaymentMinor
	}
	return 0
}
func (x *Recommendation) GetAvailability() string {
	if x != nil {
		return x.Availability
	}
	return ""
}
func (x *Recommendation) GetRank() int32 {
	if x != nil {
		return x.Rank
	}
	return 0
}

type RecommendationList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Recommendation      `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecommendationList) Reset() {
	*x = RecommendationList{}
	mi := &file_credit_v1_credit_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *RecommendationList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*RecommendationList) ProtoMessage() {
}
func (x *RecommendationList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*RecommendationList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{11}
}
func (x *RecommendationList) GetItems() []*Recommendation {
	if x != nil {
		return x.Items
	}
	return nil
}

type CalculationInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     string                 `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	AmountMinor   int64                  `protobuf:"varint,2,opt,name=amount_minor,json=amountMinor,proto3" json:"amount_minor,omitempty"`
	TermMonths    int32                  `protobuf:"varint,3,opt,name=term_months,json=termMonths,proto3" json:"term_months,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CalculationInput) Reset() {
	*x = CalculationInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CalculationInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CalculationInput) ProtoMessage() {
}
func (x *CalculationInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CalculationInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{12}
}
func (x *CalculationInput) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}
func (x *CalculationInput) GetAmountMinor() int64 {
	if x != nil {
		return x.AmountMinor
	}
	return 0
}
func (x *CalculationInput) GetTermMonths() int32 {
	if x != nil {
		return x.TermMonths
	}
	return 0
}

type Payment struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Month          int32                  `protobuf:"varint,1,opt,name=month,proto3" json:"month,omitempty"`
	PrincipalMinor int64                  `protobuf:"varint,2,opt,name=principal_minor,json=principalMinor,proto3" json:"principal_minor,omitempty"`
	InterestMinor  int64                  `protobuf:"varint,3,opt,name=interest_minor,json=interestMinor,proto3" json:"interest_minor,omitempty"`
	PaymentMinor   int64                  `protobuf:"varint,4,opt,name=payment_minor,json=paymentMinor,proto3" json:"payment_minor,omitempty"`
	BalanceMinor   int64                  `protobuf:"varint,5,opt,name=balance_minor,json=balanceMinor,proto3" json:"balance_minor,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Payment) Reset() {
	*x = Payment{}
	mi := &file_credit_v1_credit_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Payment) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Payment) ProtoMessage() {
}
func (x *Payment) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Payment) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{13}
}
func (x *Payment) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}
func (x *Payment) GetPrincipalMinor() int64 {
	if x != nil {
		return x.PrincipalMinor
	}
	return 0
}
func (x *Payment) GetInterestMinor() int64 {
	if x != nil {
		return x.InterestMinor
	}
	return 0
}
func (x *Payment) GetPaymentMinor() int64 {
	if x != nil {
		return x.PaymentMinor
	}
	return 0
}
func (x *Payment) GetBalanceMinor() int64 {
	if x != nil {
		return x.BalanceMinor
	}
	return 0
}

type Calculation struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	MonthlyPaymentMinor int64                  `protobuf:"varint,1,opt,name=monthly_payment_minor,json=monthlyPaymentMinor,proto3" json:"monthly_payment_minor,omitempty"`
	TotalPaymentMinor   int64                  `protobuf:"varint,2,opt,name=total_payment_minor,json=totalPaymentMinor,proto3" json:"total_payment_minor,omitempty"`
	OverpaymentMinor    int64                  `protobuf:"varint,3,opt,name=overpayment_minor,json=overpaymentMinor,proto3" json:"overpayment_minor,omitempty"`
	Schedule            []*Payment             `protobuf:"bytes,4,rep,name=schedule,proto3" json:"schedule,omitempty"`
	Disclaimer          string                 `protobuf:"bytes,5,opt,name=disclaimer,proto3" json:"disclaimer,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Calculation) Reset() {
	*x = Calculation{}
	mi := &file_credit_v1_credit_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Calculation) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Calculation) ProtoMessage() {
}
func (x *Calculation) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Calculation) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{14}
}
func (x *Calculation) GetMonthlyPaymentMinor() int64 {
	if x != nil {
		return x.MonthlyPaymentMinor
	}
	return 0
}
func (x *Calculation) GetTotalPaymentMinor() int64 {
	if x != nil {
		return x.TotalPaymentMinor
	}
	return 0
}
func (x *Calculation) GetOverpaymentMinor() int64 {
	if x != nil {
		return x.OverpaymentMinor
	}
	return 0
}
func (x *Calculation) GetSchedule() []*Payment {
	if x != nil {
		return x.Schedule
	}
	return nil
}
func (x *Calculation) GetDisclaimer() string {
	if x != nil {
		return x.Disclaimer
	}
	return ""
}

type Consent struct {
	state                  protoimpl.MessageState `protogen:"open.v1"`
	PersonalDataProcessing bool                   `protobuf:"varint,1,opt,name=personal_data_processing,json=personalDataProcessing,proto3" json:"personal_data_processing,omitempty"`
	CreditHistoryCheck     bool                   `protobuf:"varint,2,opt,name=credit_history_check,json=creditHistoryCheck,proto3" json:"credit_history_check,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *Consent) Reset() {
	*x = Consent{}
	mi := &file_credit_v1_credit_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Consent) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Consent) ProtoMessage() {
}
func (x *Consent) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Consent) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{15}
}
func (x *Consent) GetPersonalDataProcessing() bool {
	if x != nil {
		return x.PersonalDataProcessing
	}
	return false
}
func (x *Consent) GetCreditHistoryCheck() bool {
	if x != nil {
		return x.CreditHistoryCheck
	}
	return false
}

type DraftInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Need          *Need                  `protobuf:"bytes,1,opt,name=need,proto3" json:"need,omitempty"`
	ProductId     string                 `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Consents      *Consent               `protobuf:"bytes,3,opt,name=consents,proto3" json:"consents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DraftInput) Reset() {
	*x = DraftInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DraftInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DraftInput) ProtoMessage() {
}
func (x *DraftInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DraftInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{16}
}
func (x *DraftInput) GetNeed() *Need {
	if x != nil {
		return x.Need
	}
	return nil
}
func (x *DraftInput) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}
func (x *DraftInput) GetConsents() *Consent {
	if x != nil {
		return x.Consents
	}
	return nil
}

type DraftPatch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Need          *Need                  `protobuf:"bytes,1,opt,name=need,proto3" json:"need,omitempty"`
	ProductId     *string                `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3,oneof" json:"product_id,omitempty"`
	Consents      *Consent               `protobuf:"bytes,3,opt,name=consents,proto3" json:"consents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DraftPatch) Reset() {
	*x = DraftPatch{}
	mi := &file_credit_v1_credit_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DraftPatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DraftPatch) ProtoMessage() {
}
func (x *DraftPatch) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DraftPatch) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{17}
}
func (x *DraftPatch) GetNeed() *Need {
	if x != nil {
		return x.Need
	}
	return nil
}
func (x *DraftPatch) GetProductId() string {
	if x != nil && x.ProductId != nil {
		return *x.ProductId
	}
	return ""
}
func (x *DraftPatch) GetConsents() *Consent {
	if x != nil {
		return x.Consents
	}
	return nil
}

type Application struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CompanyId      string                 `protobuf:"bytes,2,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
	ProductId      string                 `protobuf:"bytes,3,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Need           *Need                  `protobuf:"bytes,4,opt,name=need,proto3" json:"need,omitempty"`
	Consents       *Consent               `protobuf:"bytes,5,opt,name=consents,proto3" json:"consents,omitempty"`
	Status         string                 `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	Version        int64                  `protobuf:"varint,7,opt,name=version,proto3" json:"version,omitempty"`
	CreatedAt      string                 `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      string                 `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	NextAction     string                 `protobuf:"bytes,10,opt,name=next_action,json=nextAction,proto3" json:"next_action,omitempty"`
	MissingFields  []string               `protobuf:"bytes,11,rep,name=missing_fields,json=missingFields,proto3" json:"missing_fields,omitempty"`
	Number         string                 `protobuf:"bytes,12,opt,name=number,proto3" json:"number,omitempty"`
	CompanyName    string                 `protobuf:"bytes,13,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	Inn            string                 `protobuf:"bytes,14,opt,name=inn,proto3" json:"inn,omitempty"`
	ProductName    string                 `protobuf:"bytes,15,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	Recommendation string                 `protobuf:"bytes,16,opt,name=recommendation,proto3" json:"recommendation,omitempty"`
	OwnerId        string                 `protobuf:"bytes,17,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Application) Reset() {
	*x = Application{}
	mi := &file_credit_v1_credit_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Application) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Application) ProtoMessage() {
}
func (x *Application) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Application) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{18}
}
func (x *Application) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Application) GetCompanyId() string {
	if x != nil {
		return x.CompanyId
	}
	return ""
}
func (x *Application) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}
func (x *Application) GetNeed() *Need {
	if x != nil {
		return x.Need
	}
	return nil
}
func (x *Application) GetConsents() *Consent {
	if x != nil {
		return x.Consents
	}
	return nil
}
func (x *Application) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *Application) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}
func (x *Application) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}
func (x *Application) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}
func (x *Application) GetNextAction() string {
	if x != nil {
		return x.NextAction
	}
	return ""
}
func (x *Application) GetMissingFields() []string {
	if x != nil {
		return x.MissingFields
	}
	return nil
}
func (x *Application) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}
func (x *Application) GetCompanyName() string {
	if x != nil {
		return x.CompanyName
	}
	return ""
}
func (x *Application) GetInn() string {
	if x != nil {
		return x.Inn
	}
	return ""
}
func (x *Application) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}
func (x *Application) GetRecommendation() string {
	if x != nil {
		return x.Recommendation
	}
	return ""
}
func (x *Application) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

type PageInput struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Offset         int32                  `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit          int32                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Status         string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Search         string                 `protobuf:"bytes,4,opt,name=search,proto3" json:"search,omitempty"`
	CompanyId      string                 `protobuf:"bytes,5,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
	Recommendation string                 `protobuf:"bytes,6,opt,name=recommendation,proto3" json:"recommendation,omitempty"`
	FromDate       string                 `protobuf:"bytes,7,opt,name=from_date,json=fromDate,proto3" json:"from_date,omitempty"`
	ToDate         string                 `protobuf:"bytes,8,opt,name=to_date,json=toDate,proto3" json:"to_date,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PageInput) Reset() {
	*x = PageInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *PageInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*PageInput) ProtoMessage() {
}
func (x *PageInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*PageInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{19}
}
func (x *PageInput) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}
func (x *PageInput) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}
func (x *PageInput) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *PageInput) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}
func (x *PageInput) GetCompanyId() string {
	if x != nil {
		return x.CompanyId
	}
	return ""
}
func (x *PageInput) GetRecommendation() string {
	if x != nil {
		return x.Recommendation
	}
	return ""
}
func (x *PageInput) GetFromDate() string {
	if x != nil {
		return x.FromDate
	}
	return ""
}
func (x *PageInput) GetToDate() string {
	if x != nil {
		return x.ToDate
	}
	return ""
}

type ApplicationPage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Application         `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Offset        int32                  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit         int32                  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Total         int64                  `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
	Counts        map[string]int64       `protobuf:"bytes,5,rep,name=counts,proto3" json:"counts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ApplicationPage) Reset() {
	*x = ApplicationPage{}
	mi := &file_credit_v1_credit_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ApplicationPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ApplicationPage) ProtoMessage() {
}
func (x *ApplicationPage) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ApplicationPage) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{20}
}
func (x *ApplicationPage) GetItems() []*Application {
	if x != nil {
		return x.Items
	}
	return nil
}
func (x *ApplicationPage) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}
func (x *ApplicationPage) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}
func (x *ApplicationPage) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}
func (x *ApplicationPage) GetCounts() map[string]int64 {
	if x != nil {
		return x.Counts
	}
	return nil
}

type PatchInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Patch         *DraftPatch            `protobuf:"bytes,2,opt,name=patch,proto3" json:"patch,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchInput) Reset() {
	*x = PatchInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[21]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *PatchInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*PatchInput) ProtoMessage() {
}
func (x *PatchInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[21]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*PatchInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{21}
}
func (x *PatchInput) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *PatchInput) GetPatch() *DraftPatch {
	if x != nil {
		return x.Patch
	}
	return nil
}

type Document struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ApplicationId string                 `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	Type          string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	FileName      string                 `protobuf:"bytes,4,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	ContentType   string                 `protobuf:"bytes,5,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	SizeBytes     int64                  `protobuf:"varint,6,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	Status        string                 `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
	Checksum      string                 `protobuf:"bytes,8,opt,name=checksum,proto3" json:"checksum,omitempty"`
	RequestItemId string                 `protobuf:"bytes,9,opt,name=request_item_id,json=requestItemId,proto3" json:"request_item_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Document) Reset() {
	*x = Document{}
	mi := &file_credit_v1_credit_proto_msgTypes[22]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Document) ProtoMessage() {
}
func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[22]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Document) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{22}
}
func (x *Document) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Document) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *Document) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}
func (x *Document) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}
func (x *Document) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}
func (x *Document) GetSizeBytes() int64 {
	if x != nil {
		return x.SizeBytes
	}
	return 0
}
func (x *Document) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *Document) GetChecksum() string {
	if x != nil {
		return x.Checksum
	}
	return ""
}
func (x *Document) GetRequestItemId() string {
	if x != nil {
		return x.RequestItemId
	}
	return ""
}

type DocumentList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Document            `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DocumentList) Reset() {
	*x = DocumentList{}
	mi := &file_credit_v1_credit_proto_msgTypes[23]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DocumentList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DocumentList) ProtoMessage() {
}
func (x *DocumentList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[23]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DocumentList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{23}
}
func (x *DocumentList) GetItems() []*Document {
	if x != nil {
		return x.Items
	}
	return nil
}

type UploadMetadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ApplicationId string                 `protobuf:"bytes,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	FileName      string                 `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	ContentType   string                 `protobuf:"bytes,4,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	RequestItemId string                 `protobuf:"bytes,5,opt,name=request_item_id,json=requestItemId,proto3" json:"request_item_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadMetadata) Reset() {
	*x = UploadMetadata{}
	mi := &file_credit_v1_credit_proto_msgTypes[24]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *UploadMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*UploadMetadata) ProtoMessage() {
}
func (x *UploadMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[24]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*UploadMetadata) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{24}
}
func (x *UploadMetadata) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *UploadMetadata) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}
func (x *UploadMetadata) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}
func (x *UploadMetadata) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}
func (x *UploadMetadata) GetRequestItemId() string {
	if x != nil {
		return x.RequestItemId
	}
	return ""
}

type UploadChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         isUploadChunk_Value    `protobuf_oneof:"value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadChunk) Reset() {
	*x = UploadChunk{}
	mi := &file_credit_v1_credit_proto_msgTypes[25]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *UploadChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*UploadChunk) ProtoMessage() {
}
func (x *UploadChunk) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[25]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*UploadChunk) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{25}
}
func (x *UploadChunk) GetValue() isUploadChunk_Value {
	if x != nil {
		return x.Value
	}
	return nil
}
func (x *UploadChunk) GetMetadata() *UploadMetadata {
	if x != nil {
		if x, ok := x.Value.(*UploadChunk_Metadata); ok {
			return x.Metadata
		}
	}
	return nil
}
func (x *UploadChunk) GetData() []byte {
	if x != nil {
		if x, ok := x.Value.(*UploadChunk_Data); ok {
			return x.Data
		}
	}
	return nil
}

type isUploadChunk_Value interface{ isUploadChunk_Value() }
type UploadChunk_Metadata struct {
	Metadata *UploadMetadata `protobuf:"bytes,1,opt,name=metadata,proto3,oneof"`
}
type UploadChunk_Data struct {
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

func (*UploadChunk_Metadata) isUploadChunk_Value() {
}
func (*UploadChunk_Data) isUploadChunk_Value() {
}

type DocumentID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ApplicationId string                 `protobuf:"bytes,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	DocumentId    string                 `protobuf:"bytes,2,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DocumentID) Reset() {
	*x = DocumentID{}
	mi := &file_credit_v1_credit_proto_msgTypes[26]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DocumentID) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DocumentID) ProtoMessage() {
}
func (x *DocumentID) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[26]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DocumentID) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{26}
}
func (x *DocumentID) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *DocumentID) GetDocumentId() string {
	if x != nil {
		return x.DocumentId
	}
	return ""
}

type DownloadChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileName      string                 `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	ContentType   string                 `protobuf:"bytes,2,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Data          []byte                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DownloadChunk) Reset() {
	*x = DownloadChunk{}
	mi := &file_credit_v1_credit_proto_msgTypes[27]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DownloadChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DownloadChunk) ProtoMessage() {
}
func (x *DownloadChunk) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[27]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DownloadChunk) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{27}
}
func (x *DownloadChunk) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}
func (x *DownloadChunk) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}
func (x *DownloadChunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type FieldError struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Field         string                 `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FieldError) Reset() {
	*x = FieldError{}
	mi := &file_credit_v1_credit_proto_msgTypes[28]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *FieldError) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*FieldError) ProtoMessage() {
}
func (x *FieldError) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[28]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*FieldError) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{28}
}
func (x *FieldError) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}
func (x *FieldError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Validation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Valid         bool                   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	MissingFields []string               `protobuf:"bytes,2,rep,name=missing_fields,json=missingFields,proto3" json:"missing_fields,omitempty"`
	Errors        []*FieldError          `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Validation) Reset() {
	*x = Validation{}
	mi := &file_credit_v1_credit_proto_msgTypes[29]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Validation) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Validation) ProtoMessage() {
}
func (x *Validation) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[29]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Validation) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{29}
}
func (x *Validation) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}
func (x *Validation) GetMissingFields() []string {
	if x != nil {
		return x.MissingFields
	}
	return nil
}
func (x *Validation) GetErrors() []*FieldError {
	if x != nil {
		return x.Errors
	}
	return nil
}

type Offer struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Id                  string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId           string                 `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	AmountMinor         int64                  `protobuf:"varint,3,opt,name=amount_minor,json=amountMinor,proto3" json:"amount_minor,omitempty"`
	TermMonths          int32                  `protobuf:"varint,4,opt,name=term_months,json=termMonths,proto3" json:"term_months,omitempty"`
	AnnualRateBps       int32                  `protobuf:"varint,5,opt,name=annual_rate_bps,json=annualRateBps,proto3" json:"annual_rate_bps,omitempty"`
	MonthlyPaymentMinor int64                  `protobuf:"varint,6,opt,name=monthly_payment_minor,json=monthlyPaymentMinor,proto3" json:"monthly_payment_minor,omitempty"`
	ExpiresAt           string                 `protobuf:"bytes,7,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	Collateral          string                 `protobuf:"bytes,8,opt,name=collateral,proto3" json:"collateral,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Offer) Reset() {
	*x = Offer{}
	mi := &file_credit_v1_credit_proto_msgTypes[30]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Offer) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Offer) ProtoMessage() {
}
func (x *Offer) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[30]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Offer) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{30}
}
func (x *Offer) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Offer) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}
func (x *Offer) GetAmountMinor() int64 {
	if x != nil {
		return x.AmountMinor
	}
	return 0
}
func (x *Offer) GetTermMonths() int32 {
	if x != nil {
		return x.TermMonths
	}
	return 0
}
func (x *Offer) GetAnnualRateBps() int32 {
	if x != nil {
		return x.AnnualRateBps
	}
	return 0
}
func (x *Offer) GetMonthlyPaymentMinor() int64 {
	if x != nil {
		return x.MonthlyPaymentMinor
	}
	return 0
}
func (x *Offer) GetExpiresAt() string {
	if x != nil {
		return x.ExpiresAt
	}
	return ""
}
func (x *Offer) GetCollateral() string {
	if x != nil {
		return x.Collateral
	}
	return ""
}

type Decision struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Outcome       string                 `protobuf:"bytes,2,opt,name=outcome,proto3" json:"outcome,omitempty"`
	Reason        string                 `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Offers        []*Offer               `protobuf:"bytes,4,rep,name=offers,proto3" json:"offers,omitempty"`
	DecidedAt     string                 `protobuf:"bytes,5,opt,name=decided_at,json=decidedAt,proto3" json:"decided_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Decision) Reset() {
	*x = Decision{}
	mi := &file_credit_v1_credit_proto_msgTypes[31]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Decision) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Decision) ProtoMessage() {
}
func (x *Decision) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[31]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Decision) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{31}
}
func (x *Decision) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Decision) GetOutcome() string {
	if x != nil {
		return x.Outcome
	}
	return ""
}
func (x *Decision) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}
func (x *Decision) GetOffers() []*Offer {
	if x != nil {
		return x.Offers
	}
	return nil
}
func (x *Decision) GetDecidedAt() string {
	if x != nil {
		return x.DecidedAt
	}
	return ""
}

type InformationItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Type          string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Required      bool                   `protobuf:"varint,4,opt,name=required,proto3" json:"required,omitempty"`
	Status        string                 `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InformationItem) Reset() {
	*x = InformationItem{}
	mi := &file_credit_v1_credit_proto_msgTypes[32]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *InformationItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*InformationItem) ProtoMessage() {
}
func (x *InformationItem) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[32]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*InformationItem) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{32}
}
func (x *InformationItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *InformationItem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}
func (x *InformationItem) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}
func (x *InformationItem) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}
func (x *InformationItem) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type InformationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Items         []*InformationItem     `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InformationRequest) Reset() {
	*x = InformationRequest{}
	mi := &file_credit_v1_credit_proto_msgTypes[33]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *InformationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*InformationRequest) ProtoMessage() {
}
func (x *InformationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[33]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*InformationRequest) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{33}
}
func (x *InformationRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *InformationRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}
func (x *InformationRequest) GetItems() []*InformationItem {
	if x != nil {
		return x.Items
	}
	return nil
}
func (x *InformationRequest) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type InformationList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*InformationRequest  `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InformationList) Reset() {
	*x = InformationList{}
	mi := &file_credit_v1_credit_proto_msgTypes[34]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *InformationList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*InformationList) ProtoMessage() {
}
func (x *InformationList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[34]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*InformationList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{34}
}
func (x *InformationList) GetItems() []*InformationRequest {
	if x != nil {
		return x.Items
	}
	return nil
}

type ManualDecision struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Outcome         string                 `protobuf:"bytes,1,opt,name=outcome,proto3" json:"outcome,omitempty"`
	Reason          string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	Offers          []*Offer               `protobuf:"bytes,3,rep,name=offers,proto3" json:"offers,omitempty"`
	RequestedFields []string               `protobuf:"bytes,4,rep,name=requested_fields,json=requestedFields,proto3" json:"requested_fields,omitempty"`
	Items           []*InformationItem     `protobuf:"bytes,5,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ManualDecision) Reset() {
	*x = ManualDecision{}
	mi := &file_credit_v1_credit_proto_msgTypes[35]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ManualDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ManualDecision) ProtoMessage() {
}
func (x *ManualDecision) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[35]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ManualDecision) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{35}
}
func (x *ManualDecision) GetOutcome() string {
	if x != nil {
		return x.Outcome
	}
	return ""
}
func (x *ManualDecision) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}
func (x *ManualDecision) GetOffers() []*Offer {
	if x != nil {
		return x.Offers
	}
	return nil
}
func (x *ManualDecision) GetRequestedFields() []string {
	if x != nil {
		return x.RequestedFields
	}
	return nil
}
func (x *ManualDecision) GetItems() []*InformationItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ManualInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Decision      *ManualDecision        `protobuf:"bytes,2,opt,name=decision,proto3" json:"decision,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ManualInput) Reset() {
	*x = ManualInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[36]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ManualInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ManualInput) ProtoMessage() {
}
func (x *ManualInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[36]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ManualInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{36}
}
func (x *ManualInput) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *ManualInput) GetDecision() *ManualDecision {
	if x != nil {
		return x.Decision
	}
	return nil
}

type DeclineInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reason        string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeclineInput) Reset() {
	*x = DeclineInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[37]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *DeclineInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*DeclineInput) ProtoMessage() {
}
func (x *DeclineInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[37]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*DeclineInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{37}
}
func (x *DeclineInput) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *DeclineInput) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type AcceptInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OfferId       string                 `protobuf:"bytes,2,opt,name=offer_id,json=offerId,proto3" json:"offer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AcceptInput) Reset() {
	*x = AcceptInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[38]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *AcceptInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*AcceptInput) ProtoMessage() {
}
func (x *AcceptInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[38]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*AcceptInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{38}
}
func (x *AcceptInput) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *AcceptInput) GetOfferId() string {
	if x != nil {
		return x.OfferId
	}
	return ""
}

type Loan struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ApplicationId string                 `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	OfferId       string                 `protobuf:"bytes,3,opt,name=offer_id,json=offerId,proto3" json:"offer_id,omitempty"`
	AmountMinor   int64                  `protobuf:"varint,4,opt,name=amount_minor,json=amountMinor,proto3" json:"amount_minor,omitempty"`
	IssuedAt      string                 `protobuf:"bytes,5,opt,name=issued_at,json=issuedAt,proto3" json:"issued_at,omitempty"`
	Simulated     bool                   `protobuf:"varint,6,opt,name=simulated,proto3" json:"simulated,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Loan) Reset() {
	*x = Loan{}
	mi := &file_credit_v1_credit_proto_msgTypes[39]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Loan) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Loan) ProtoMessage() {
}
func (x *Loan) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[39]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Loan) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{39}
}
func (x *Loan) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Loan) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *Loan) GetOfferId() string {
	if x != nil {
		return x.OfferId
	}
	return ""
}
func (x *Loan) GetAmountMinor() int64 {
	if x != nil {
		return x.AmountMinor
	}
	return 0
}
func (x *Loan) GetIssuedAt() string {
	if x != nil {
		return x.IssuedAt
	}
	return ""
}
func (x *Loan) GetSimulated() bool {
	if x != nil {
		return x.Simulated
	}
	return false
}

type Event struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	At            string                 `protobuf:"bytes,2,opt,name=at,proto3" json:"at,omitempty"`
	FromStatus    string                 `protobuf:"bytes,3,opt,name=from_status,json=fromStatus,proto3" json:"from_status,omitempty"`
	ToStatus      string                 `protobuf:"bytes,4,opt,name=to_status,json=toStatus,proto3" json:"to_status,omitempty"`
	Message       string                 `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
	ActorRole     string                 `protobuf:"bytes,6,opt,name=actor_role,json=actorRole,proto3" json:"actor_role,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Event) Reset() {
	*x = Event{}
	mi := &file_credit_v1_credit_proto_msgTypes[40]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Event) ProtoMessage() {
}
func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[40]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Event) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{40}
}
func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *Event) GetAt() string {
	if x != nil {
		return x.At
	}
	return ""
}
func (x *Event) GetFromStatus() string {
	if x != nil {
		return x.FromStatus
	}
	return ""
}
func (x *Event) GetToStatus() string {
	if x != nil {
		return x.ToStatus
	}
	return ""
}
func (x *Event) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}
func (x *Event) GetActorRole() string {
	if x != nil {
		return x.ActorRole
	}
	return ""
}

type EventList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Event               `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventList) Reset() {
	*x = EventList{}
	mi := &file_credit_v1_credit_proto_msgTypes[41]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *EventList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*EventList) ProtoMessage() {
}
func (x *EventList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[41]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*EventList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{41}
}
func (x *EventList) GetItems() []*Event {
	if x != nil {
		return x.Items
	}
	return nil
}

type Check struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Source        string                 `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Summary       string                 `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	CheckedAt     string                 `protobuf:"bytes,4,opt,name=checked_at,json=checkedAt,proto3" json:"checked_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Check) Reset() {
	*x = Check{}
	mi := &file_credit_v1_credit_proto_msgTypes[42]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Check) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Check) ProtoMessage() {
}
func (x *Check) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[42]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Check) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{42}
}
func (x *Check) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}
func (x *Check) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *Check) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}
func (x *Check) GetCheckedAt() string {
	if x != nil {
		return x.CheckedAt
	}
	return ""
}

type Affordability struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Status          string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	CfadsMinor      int64                  `protobuf:"varint,2,opt,name=cfads_minor,json=cfadsMinor,proto3" json:"cfads_minor,omitempty"`
	MaxPaymentMinor int64                  `protobuf:"varint,3,opt,name=max_payment_minor,json=maxPaymentMinor,proto3" json:"max_payment_minor,omitempty"`
	MaxAmountMinor  int64                  `protobuf:"varint,4,opt,name=max_amount_minor,json=maxAmountMinor,proto3" json:"max_amount_minor,omitempty"`
	Dscr            float64                `protobuf:"fixed64,5,opt,name=dscr,proto3" json:"dscr,omitempty"`
	Reason          string                 `protobuf:"bytes,6,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Affordability) Reset() {
	*x = Affordability{}
	mi := &file_credit_v1_credit_proto_msgTypes[43]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Affordability) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Affordability) ProtoMessage() {
}
func (x *Affordability) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[43]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Affordability) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{43}
}
func (x *Affordability) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *Affordability) GetCfadsMinor() int64 {
	if x != nil {
		return x.CfadsMinor
	}
	return 0
}
func (x *Affordability) GetMaxPaymentMinor() int64 {
	if x != nil {
		return x.MaxPaymentMinor
	}
	return 0
}
func (x *Affordability) GetMaxAmountMinor() int64 {
	if x != nil {
		return x.MaxAmountMinor
	}
	return 0
}
func (x *Affordability) GetDscr() float64 {
	if x != nil {
		return x.Dscr
	}
	return 0
}
func (x *Affordability) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type Scoring struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Score          int32                  `protobuf:"varint,1,opt,name=score,proto3" json:"score,omitempty"`
	RiskLevel      string                 `protobuf:"bytes,2,opt,name=risk_level,json=riskLevel,proto3" json:"risk_level,omitempty"`
	Factors        []string               `protobuf:"bytes,3,rep,name=factors,proto3" json:"factors,omitempty"`
	ModelVersion   string                 `protobuf:"bytes,4,opt,name=model_version,json=modelVersion,proto3" json:"model_version,omitempty"`
	Checks         []*Check               `protobuf:"bytes,5,rep,name=checks,proto3" json:"checks,omitempty"`
	Recommendation string                 `protobuf:"bytes,6,opt,name=recommendation,proto3" json:"recommendation,omitempty"`
	Affordability  *Affordability         `protobuf:"bytes,7,opt,name=affordability,proto3" json:"affordability,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Scoring) Reset() {
	*x = Scoring{}
	mi := &file_credit_v1_credit_proto_msgTypes[44]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Scoring) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Scoring) ProtoMessage() {
}
func (x *Scoring) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[44]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Scoring) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{44}
}
func (x *Scoring) GetScore() int32 {
	if x != nil {
		return x.Score
	}
	return 0
}
func (x *Scoring) GetRiskLevel() string {
	if x != nil {
		return x.RiskLevel
	}
	return ""
}
func (x *Scoring) GetFactors() []string {
	if x != nil {
		return x.Factors
	}
	return nil
}
func (x *Scoring) GetModelVersion() string {
	if x != nil {
		return x.ModelVersion
	}
	return ""
}
func (x *Scoring) GetChecks() []*Check {
	if x != nil {
		return x.Checks
	}
	return nil
}
func (x *Scoring) GetRecommendation() string {
	if x != nil {
		return x.Recommendation
	}
	return ""
}
func (x *Scoring) GetAffordability() *Affordability {
	if x != nil {
		return x.Affordability
	}
	return nil
}

type RegistryData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	LegalForm     string                 `protobuf:"bytes,2,opt,name=legal_form,json=legalForm,proto3" json:"legal_form,omitempty"`
	OkvedSection  string                 `protobuf:"bytes,3,opt,name=okved_section,json=okvedSection,proto3" json:"okved_section,omitempty"`
	RegisteredAt  string                 `protobuf:"bytes,4,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at,omitempty"`
	TaxDebtMinor  int64                  `protobuf:"varint,5,opt,name=tax_debt_minor,json=taxDebtMinor,proto3" json:"tax_debt_minor,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegistryData) Reset() {
	*x = RegistryData{}
	mi := &file_credit_v1_credit_proto_msgTypes[45]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *RegistryData) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*RegistryData) ProtoMessage() {
}
func (x *RegistryData) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[45]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*RegistryData) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{45}
}
func (x *RegistryData) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}
func (x *RegistryData) GetLegalForm() string {
	if x != nil {
		return x.LegalForm
	}
	return ""
}
func (x *RegistryData) GetOkvedSection() string {
	if x != nil {
		return x.OkvedSection
	}
	return ""
}
func (x *RegistryData) GetRegisteredAt() string {
	if x != nil {
		return x.RegisteredAt
	}
	return ""
}
func (x *RegistryData) GetTaxDebtMinor() int64 {
	if x != nil {
		return x.TaxDebtMinor
	}
	return 0
}

type CreditHistory struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Found              bool                   `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	OverdueCount_12M   int32                  `protobuf:"varint,2,opt,name=overdue_count_12m,json=overdueCount12m,proto3" json:"overdue_count_12m,omitempty"`
	MaxOverdueDays_24M int32                  `protobuf:"varint,3,opt,name=max_overdue_days_24m,json=maxOverdueDays24m,proto3" json:"max_overdue_days_24m,omitempty"`
	ActiveLoans        int32                  `protobuf:"varint,4,opt,name=active_loans,json=activeLoans,proto3" json:"active_loans,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *CreditHistory) Reset() {
	*x = CreditHistory{}
	mi := &file_credit_v1_credit_proto_msgTypes[46]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CreditHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CreditHistory) ProtoMessage() {
}
func (x *CreditHistory) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[46]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CreditHistory) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{46}
}
func (x *CreditHistory) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}
func (x *CreditHistory) GetOverdueCount_12M() int32 {
	if x != nil {
		return x.OverdueCount_12M
	}
	return 0
}
func (x *CreditHistory) GetMaxOverdueDays_24M() int32 {
	if x != nil {
		return x.MaxOverdueDays_24M
	}
	return 0
}
func (x *CreditHistory) GetActiveLoans() int32 {
	if x != nil {
		return x.ActiveLoans
	}
	return 0
}

type CreditReport struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Company       *CreditHistory         `protobuf:"bytes,1,opt,name=company,proto3" json:"company,omitempty"`
	Owner         *CreditHistory         `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreditReport) Reset() {
	*x = CreditReport{}
	mi := &file_credit_v1_credit_proto_msgTypes[47]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *CreditReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*CreditReport) ProtoMessage() {
}
func (x *CreditReport) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[47]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*CreditReport) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{47}
}
func (x *CreditReport) GetCompany() *CreditHistory {
	if x != nil {
		return x.Company
	}
	return nil
}
func (x *CreditReport) GetOwner() *CreditHistory {
	if x != nil {
		return x.Owner
	}
	return nil
}

type AccountData struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	Found                bool                   `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	TurnoverMonthlyMinor int64                  `protobuf:"varint,2,opt,name=turnover_monthly_minor,json=turnoverMonthlyMinor,proto3" json:"turnover_monthly_minor,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *AccountData) Reset() {
	*x = AccountData{}
	mi := &file_credit_v1_credit_proto_msgTypes[48]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *AccountData) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*AccountData) ProtoMessage() {
}
func (x *AccountData) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[48]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*AccountData) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{48}
}
func (x *AccountData) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}
func (x *AccountData) GetTurnoverMonthlyMinor() int64 {
	if x != nil {
		return x.TurnoverMonthlyMinor
	}
	return 0
}

type Features struct {
	state                    protoimpl.MessageState `protogen:"open.v1"`
	BusinessAgeMonths        int32                  `protobuf:"varint,1,opt,name=business_age_months,json=businessAgeMonths,proto3" json:"business_age_months,omitempty"`
	MonthlyRevenueMinor      int64                  `protobuf:"varint,2,opt,name=monthly_revenue_minor,json=monthlyRevenueMinor,proto3" json:"monthly_revenue_minor,omitempty"`
	MonthlyExpensesMinor     int64                  `protobuf:"varint,3,opt,name=monthly_expenses_minor,json=monthlyExpensesMinor,proto3" json:"monthly_expenses_minor,omitempty"`
	MonthlyDebtPaymentsMinor int64                  `protobuf:"varint,4,opt,name=monthly_debt_payments_minor,json=monthlyDebtPaymentsMinor,proto3" json:"monthly_debt_payments_minor,omitempty"`
	RequestedAmountMinor     int64                  `protobuf:"varint,5,opt,name=requested_amount_minor,json=requestedAmountMinor,proto3" json:"requested_amount_minor,omitempty"`
	TermMonths               int32                  `protobuf:"varint,6,opt,name=term_months,json=termMonths,proto3" json:"term_months,omitempty"`
	Registry                 *RegistryData          `protobuf:"bytes,8,opt,name=registry,proto3" json:"registry,omitempty"`
	CreditReport             *CreditReport          `protobuf:"bytes,9,opt,name=credit_report,json=creditReport,proto3" json:"credit_report,omitempty"`
	Account                  *AccountData           `protobuf:"bytes,10,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *Features) Reset() {
	*x = Features{}
	mi := &file_credit_v1_credit_proto_msgTypes[49]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Features) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Features) ProtoMessage() {
}
func (x *Features) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[49]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Features) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{49}
}
func (x *Features) GetBusinessAgeMonths() int32 {
	if x != nil {
		return x.BusinessAgeMonths
	}
	return 0
}
func (x *Features) GetMonthlyRevenueMinor() int64 {
	if x != nil {
		return x.MonthlyRevenueMinor
	}
	return 0
}
func (x *Features) GetMonthlyExpensesMinor() int64 {
	if x != nil {
		return x.MonthlyExpensesMinor
	}
	return 0
}
func (x *Features) GetMonthlyDebtPaymentsMinor() int64 {
	if x != nil {
		return x.MonthlyDebtPaymentsMinor
	}
	return 0
}
func (x *Features) GetRequestedAmountMinor() int64 {
	if x != nil {
		return x.RequestedAmountMinor
	}
	return 0
}
func (x *Features) GetTermMonths() int32 {
	if x != nil {
		return x.TermMonths
	}
	return 0
}
func (x *Features) GetRegistry() *RegistryData {
	if x != nil {
		return x.Registry
	}
	return nil
}
func (x *Features) GetCreditReport() *CreditReport {
	if x != nil {
		return x.CreditReport
	}
	return nil
}
func (x *Features) GetAccount() *AccountData {
	if x != nil {
		return x.Account
	}
	return nil
}

type ScoringInput struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	ScoringRequestId   string                 `protobuf:"bytes,1,opt,name=scoring_request_id,json=scoringRequestId,proto3" json:"scoring_request_id,omitempty"`
	ApplicationId      string                 `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	ApplicationVersion int64                  `protobuf:"varint,3,opt,name=application_version,json=applicationVersion,proto3" json:"application_version,omitempty"`
	Features           *Features              `protobuf:"bytes,4,opt,name=features,proto3" json:"features,omitempty"`
	Product            *Product               `protobuf:"bytes,5,opt,name=product,proto3" json:"product,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *ScoringInput) Reset() {
	*x = ScoringInput{}
	mi := &file_credit_v1_credit_proto_msgTypes[50]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ScoringInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ScoringInput) ProtoMessage() {
}
func (x *ScoringInput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[50]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ScoringInput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{50}
}
func (x *ScoringInput) GetScoringRequestId() string {
	if x != nil {
		return x.ScoringRequestId
	}
	return ""
}
func (x *ScoringInput) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *ScoringInput) GetApplicationVersion() int64 {
	if x != nil {
		return x.ApplicationVersion
	}
	return 0
}
func (x *ScoringInput) GetFeatures() *Features {
	if x != nil {
		return x.Features
	}
	return nil
}
func (x *ScoringInput) GetProduct() *Product {
	if x != nil {
		return x.Product
	}
	return nil
}

type ScoringOutput struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	ScoringRequestId   string                 `protobuf:"bytes,1,opt,name=scoring_request_id,json=scoringRequestId,proto3" json:"scoring_request_id,omitempty"`
	ApplicationId      string                 `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	ApplicationVersion int64                  `protobuf:"varint,3,opt,name=application_version,json=applicationVersion,proto3" json:"application_version,omitempty"`
	Scoring            *Scoring               `protobuf:"bytes,4,opt,name=scoring,proto3" json:"scoring,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *ScoringOutput) Reset() {
	*x = ScoringOutput{}
	mi := &file_credit_v1_credit_proto_msgTypes[51]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *ScoringOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ScoringOutput) ProtoMessage() {
}
func (x *ScoringOutput) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[51]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ScoringOutput) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{51}
}
func (x *ScoringOutput) GetScoringRequestId() string {
	if x != nil {
		return x.ScoringRequestId
	}
	return ""
}
func (x *ScoringOutput) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}
func (x *ScoringOutput) GetApplicationVersion() int64 {
	if x != nil {
		return x.ApplicationVersion
	}
	return 0
}
func (x *ScoringOutput) GetScoring() *Scoring {
	if x != nil {
		return x.Scoring
	}
	return nil
}

type FeedItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Body          string                 `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Reason        string                 `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ApplicationId string                 `protobuf:"bytes,6,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FeedItem) Reset() {
	*x = FeedItem{}
	mi := &file_credit_v1_credit_proto_msgTypes[52]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *FeedItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*FeedItem) ProtoMessage() {
}
func (x *FeedItem) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[52]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*FeedItem) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{52}
}
func (x *FeedItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}
func (x *FeedItem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}
func (x *FeedItem) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}
func (x *FeedItem) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}
func (x *FeedItem) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}
func (x *FeedItem) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}

type FeedList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*FeedItem            `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FeedList) Reset() {
	*x = FeedList{}
	mi := &file_credit_v1_credit_proto_msgTypes[53]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *FeedList) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*FeedList) ProtoMessage() {
}
func (x *FeedList) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[53]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*FeedList) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{53}
}
func (x *FeedList) GetItems() []*FeedItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type Overview struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Application         *Application           `protobuf:"bytes,1,opt,name=application,proto3" json:"application,omitempty"`
	CompanySnapshot     *Company               `protobuf:"bytes,2,opt,name=company_snapshot,json=companySnapshot,proto3" json:"company_snapshot,omitempty"`
	ProductSnapshot     *Product               `protobuf:"bytes,3,opt,name=product_snapshot,json=productSnapshot,proto3" json:"product_snapshot,omitempty"`
	Scoring             *Scoring               `protobuf:"bytes,4,opt,name=scoring,proto3" json:"scoring,omitempty"`
	Decision            *Decision              `protobuf:"bytes,5,opt,name=decision,proto3" json:"decision,omitempty"`
	Documents           []*Document            `protobuf:"bytes,6,rep,name=documents,proto3" json:"documents,omitempty"`
	InformationRequests []*InformationRequest  `protobuf:"bytes,7,rep,name=information_requests,json=informationRequests,proto3" json:"information_requests,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Overview) Reset() {
	*x = Overview{}
	mi := &file_credit_v1_credit_proto_msgTypes[54]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}
func (x *Overview) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Overview) ProtoMessage() {
}
func (x *Overview) ProtoReflect() protoreflect.Message {
	mi := &file_credit_v1_credit_proto_msgTypes[54]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Overview) Descriptor() ([]byte, []int) {
	return file_credit_v1_credit_proto_rawDescGZIP(), []int{54}
}
func (x *Overview) GetApplication() *Application {
	if x != nil {
		return x.Application
	}
	return nil
}
func (x *Overview) GetCompanySnapshot() *Company {
	if x != nil {
		return x.CompanySnapshot
	}
	return nil
}
func (x *Overview) GetProductSnapshot() *Product {
	if x != nil {
		return x.ProductSnapshot
	}
	return nil
}
func (x *Overview) GetScoring() *Scoring {
	if x != nil {
		return x.Scoring
	}
	return nil
}
func (x *Overview) GetDecision() *Decision {
	if x != nil {
		return x.Decision
	}
	return nil
}
func (x *Overview) GetDocuments() []*Document {
	if x != nil {
		return x.Documents
	}
	return nil
}
func (x *Overview) GetInformationRequests() []*InformationRequest {
	if x != nil {
		return x.InformationRequests
	}
	return nil
}

var File_credit_v1_credit_proto protoreflect.FileDescriptor

const file_credit_v1_credit_proto_rawDesc = "" + "\n" + "\x16credit/v1/credit.proto\x12\tcredit.v1\"\a\n" + "\x05Empty\"\x14\n" + "\x02ID\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\"\x81\x06\n" + "\fCompanyInput\x12\x10\n" + "\x03inn\x18\x01 \x01(\tR\x03inn\x12\x12\n" + "\x04name\x18\x02 \x01(\tR\x04name\x12.\n" + "\x13business_age_months\x18\x03 \x01(\x05R\x11businessAgeMonths\x122\n" + "\x15monthly_revenue_minor\x18\x04 \x01(\x03R\x13monthlyRevenueMinor\x124\n" + "\x16monthly_expenses_minor\x18\x05 \x01(\x03R\x14monthlyExpensesMinor\x12=\n" + "\x1bmonthly_debt_payments_minor\x18\x06 \x01(\x03R\x18monthlyDebtPaymentsMinor\x12\x10\n" + "\x03kpp\x18\a \x01(\tR\x03kpp\x12\x12\n" + "\x04ogrn\x18\b \x01(\tR\x04ogrn\x12\x1a\n" + "\bindustry\x18\t \x01(\tR\bindustry\x12#\n" + "\rregistered_at\x18\n" + " \x01(\tR\fregisteredAt\x12/\n" + "\x13representative_name\x18\v \x01(\tR\x12representativeName\x127\n" + "\x17representative_position\x18\f \x01(\tR\x16representativePosition\x12\x14\n" + "\x05phone\x18\r \x01(\tR\x05phone\x12\x14\n" + "\x05email\x18\x0e \x01(\tR\x05email\x120\n" + "\x14annual_revenue_minor\x18\x0f \x01(\x03R\x12annualRevenueMinor\x122\n" + "\x15annual_expenses_minor\x18\x10 \x01(\x03R\x13annualExpensesMinor\x12&\n" + "\x0fcash_flow_minor\x18\x11 \x01(\x03R\rcashFlowMinor\x124\n" + "\x16outstanding_debt_minor\x18\x12 \x01(\x03R\x14outstandingDebtMinor\x12+\n" + "\x11marketplace_sales\x18\x14 \x01(\bR\x10marketplaceSalesJ\x04\b\x13\x10\x14\"\xdf\x01\n" + "\aCompany\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12+\n" + "\x04data\x18\x02 \x01(\v2\x17.credit.v1.CompanyInputR\x04data\x12/\n" + "\x13verification_status\x18\x03 \x01(\tR\x12verificationStatus\x121\n" + "\x14verification_message\x18\x04 \x01(\tR\x13verificationMessage\x12\x19\n" + "\bowner_id\x18\x05 \x01(\tR\aownerId\x12\x18\n" + "\aversion\x18\x06 \x01(\x03R\aversion\"7\n" + "\vCompanyList\x12(\n" + "\x05items\x18\x01 \x03(\v2\x12.credit.v1.CompanyR\x05items\"L\n" + "\rCompanyUpdate\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12+\n" + "\x04data\x18\x02 \x01(\v2\x17.credit.v1.CompanyInputR\x04data\"\xc6\x03\n" + "\aProduct\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" + "\x04name\x18\x02 \x01(\tR\x04name\x12\x1a\n" + "\bpurposes\x18\x03 \x03(\tR\bpurposes\x12(\n" + "\x10min_amount_minor\x18\x04 \x01(\x03R\x0eminAmountMinor\x12(\n" + "\x10max_amount_minor\x18\x05 \x01(\x03R\x0emaxAmountMinor\x12&\n" + "\x0fmin_term_months\x18\x06 \x01(\x05R\rminTermMonths\x12&\n" + "\x0fmax_term_months\x18\a \x01(\x05R\rmaxTermMonths\x12&\n" + "\x0fannual_rate_bps\x18\b \x01(\x05R\rannualRateBps\x126\n" + "\x17required_document_types\x18\t \x03(\tR\x15requiredDocumentTypes\x12 \n" + "\vdescription\x18\n" + " \x01(\tR\vdescription\x12\x1e\n" + "\n" + "collateral\x18\v \x01(\tR\n" + "collateral\x12!\n" + "\fpayment_type\x18\f \x01(\tR\vpaymentType\x12\x12\n" + "\x04tags\x18\r \x03(\tR\x04tags\"7\n" + "\vProductList\x12(\n" + "\x05items\x18\x01 \x03(\v2\x12.credit.v1.ProductR\x05items\"\xa5\x01\n" + "\x04Need\x12\x1d\n" + "\n" + "company_id\x18\x01 \x01(\tR\tcompanyId\x12!\n" + "\famount_minor\x18\x02 \x01(\x03R\vamountMinor\x12\x1f\n" + "\vterm_months\x18\x03 \x01(\x05R\n" + "termMonths\x12\x18\n" + "\apurpose\x18\x04 \x01(\tR\apurpose\x12 \n" + "\vdescription\x18\x05 \x01(\tR\vdescription\"\xc8\x01\n" + "\x13RecommendationInput\x12\x1d\n" + "\n" + "company_id\x18\x01 \x01(\tR\tcompanyId\x12&\n" + "\famount_minor\x18\x02 \x01(\x03H\x00R\vamountMinor\x88\x01\x01\x12\x1f\n" + "\vterm_months\x18\x03 \x01(\x05R\n" + "termMonths\x12\x18\n" + "\apurpose\x18\x04 \x01(\tR\apurpose\x12\x1e\n" + "\n" + "collateral\x18\x05 \x01(\tR\n" + "collateralB\x0f\n" + "\r_amount_minor\"\xd5\x01\n" + "\x0eRecommendation\x12,\n" + "\aproduct\x18\x01 \x01(\v2\x12.credit.v1.ProductR\aproduct\x12\x16\n" + "\x06reason\x18\x02 \x01(\tR\x06reason\x12E\n" + "\x1festimated_monthly_payment_minor\x18\x03 \x01(\x03R\x1cestimatedMonthlyPaymentMinor\x12\"\n" + "\favailability\x18\x04 \x01(\tR\favailability\x12\x12\n" + "\x04rank\x18\x05 \x01(\x05R\x04rank\"E\n" + "\x12RecommendationList\x12/\n" + "\x05items\x18\x01 \x03(\v2\x19.credit.v1.RecommendationR\x05items\"u\n" + "\x10CalculationInput\x12\x1d\n" + "\n" + "product_id\x18\x01 \x01(\tR\tproductId\x12!\n" + "\famount_minor\x18\x02 \x01(\x03R\vamountMinor\x12\x1f\n" + "\vterm_months\x18\x03 \x01(\x05R\n" + "termMonths\"\xb9\x01\n" + "\aPayment\x12\x14\n" + "\x05month\x18\x01 \x01(\x05R\x05month\x12'\n" + "\x0fprincipal_minor\x18\x02 \x01(\x03R\x0eprincipalMinor\x12%\n" + "\x0einterest_minor\x18\x03 \x01(\x03R\rinterestMinor\x12#\n" + "\rpayment_minor\x18\x04 \x01(\x03R\fpaymentMinor\x12#\n" + "\rbalance_minor\x18\x05 \x01(\x03R\fbalanceMinor\"\xee\x01\n" + "\vCalculation\x122\n" + "\x15monthly_payment_minor\x18\x01 \x01(\x03R\x13monthlyPaymentMinor\x12.\n" + "\x13total_payment_minor\x18\x02 \x01(\x03R\x11totalPaymentMinor\x12+\n" + "\x11overpayment_minor\x18\x03 \x01(\x03R\x10overpaymentMinor\x12.\n" + "\bschedule\x18\x04 \x03(\v2\x12.credit.v1.PaymentR\bschedule\x12\x1e\n" + "\n" + "disclaimer\x18\x05 \x01(\tR\n" + "disclaimer\"u\n" + "\aConsent\x128\n" + "\x18personal_data_processing\x18\x01 \x01(\bR\x16personalDataProcessing\x120\n" + "\x14credit_history_check\x18\x02 \x01(\bR\x12creditHistoryCheck\"\x80\x01\n" + "\n" + "DraftInput\x12#\n" + "\x04need\x18\x01 \x01(\v2\x0f.credit.v1.NeedR\x04need\x12\x1d\n" + "\n" + "product_id\x18\x02 \x01(\tR\tproductId\x12.\n" + "\bconsents\x18\x03 \x01(\v2\x12.credit.v1.ConsentR\bconsents\"\x94\x01\n" + "\n" + "DraftPatch\x12#\n" + "\x04need\x18\x01 \x01(\v2\x0f.credit.v1.NeedR\x04need\x12\"\n" + "\n" + "product_id\x18\x02 \x01(\tH\x00R\tproductId\x88\x01\x01\x12.\n" + "\bconsents\x18\x03 \x01(\v2\x12.credit.v1.ConsentR\bconsentsB\r\n" + "\v_product_id\"\x9b\x04\n" + "\vApplication\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x1d\n" + "\n" + "company_id\x18\x02 \x01(\tR\tcompanyId\x12\x1d\n" + "\n" + "product_id\x18\x03 \x01(\tR\tproductId\x12#\n" + "\x04need\x18\x04 \x01(\v2\x0f.credit.v1.NeedR\x04need\x12.\n" + "\bconsents\x18\x05 \x01(\v2\x12.credit.v1.ConsentR\bconsents\x12\x16\n" + "\x06status\x18\x06 \x01(\tR\x06status\x12\x18\n" + "\aversion\x18\a \x01(\x03R\aversion\x12\x1d\n" + "\n" + "created_at\x18\b \x01(\tR\tcreatedAt\x12\x1d\n" + "\n" + "updated_at\x18\t \x01(\tR\tupdatedAt\x12\x1f\n" + "\vnext_action\x18\n" + " \x01(\tR\n" + "nextAction\x12%\n" + "\x0emissing_fields\x18\v \x03(\tR\rmissingFields\x12\x16\n" + "\x06number\x18\f \x01(\tR\x06number\x12!\n" + "\fcompany_name\x18\r \x01(\tR\vcompanyName\x12\x10\n" + "\x03inn\x18\x0e \x01(\tR\x03inn\x12!\n" + "\fproduct_name\x18\x0f \x01(\tR\vproductName\x12&\n" + "\x0erecommendation\x18\x10 \x01(\tR\x0erecommendation\x12\x19\n" + "\bowner_id\x18\x11 \x01(\tR\aownerId\"\xe6\x01\n" + "\tPageInput\x12\x16\n" + "\x06offset\x18\x01 \x01(\x05R\x06offset\x12\x14\n" + "\x05limit\x18\x02 \x01(\x05R\x05limit\x12\x16\n" + "\x06status\x18\x03 \x01(\tR\x06status\x12\x16\n" + "\x06search\x18\x04 \x01(\tR\x06search\x12\x1d\n" + "\n" + "company_id\x18\x05 \x01(\tR\tcompanyId\x12&\n" + "\x0erecommendation\x18\x06 \x01(\tR\x0erecommendation\x12\x1b\n" + "\tfrom_date\x18\a \x01(\tR\bfromDate\x12\x17\n" + "\ato_date\x18\b \x01(\tR\x06toDate\"\xfe\x01\n" + "\x0fApplicationPage\x12,\n" + "\x05items\x18\x01 \x03(\v2\x16.credit.v1.ApplicationR\x05items\x12\x16\n" + "\x06offset\x18\x02 \x01(\x05R\x06offset\x12\x14\n" + "\x05limit\x18\x03 \x01(\x05R\x05limit\x12\x14\n" + "\x05total\x18\x04 \x01(\x03R\x05total\x12>\n" + "\x06counts\x18\x05 \x03(\v2&.credit.v1.ApplicationPage.CountsEntryR\x06counts\x1a9\n" + "\vCountsEntry\x12\x10\n" + "\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" + "\x05value\x18\x02 \x01(\x03R\x05value:\x028\x01\"I\n" + "\n" + "PatchInput\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12+\n" + "\x05patch\x18\x02 \x01(\v2\x15.credit.v1.DraftPatchR\x05patch\"\x90\x02\n" + "\bDocument\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12%\n" + "\x0eapplication_id\x18\x02 \x01(\tR\rapplicationId\x12\x12\n" + "\x04type\x18\x03 \x01(\tR\x04type\x12\x1b\n" + "\tfile_name\x18\x04 \x01(\tR\bfileName\x12!\n" + "\fcontent_type\x18\x05 \x01(\tR\vcontentType\x12\x1d\n" + "\n" + "size_bytes\x18\x06 \x01(\x03R\tsizeBytes\x12\x16\n" + "\x06status\x18\a \x01(\tR\x06status\x12\x1a\n" + "\bchecksum\x18\b \x01(\tR\bchecksum\x12&\n" + "\x0frequest_item_id\x18\t \x01(\tR\rrequestItemId\"9\n" + "\fDocumentList\x12)\n" + "\x05items\x18\x01 \x03(\v2\x13.credit.v1.DocumentR\x05items\"\xb3\x01\n" + "\x0eUploadMetadata\x12%\n" + "\x0eapplication_id\x18\x01 \x01(\tR\rapplicationId\x12\x12\n" + "\x04type\x18\x02 \x01(\tR\x04type\x12\x1b\n" + "\tfile_name\x18\x03 \x01(\tR\bfileName\x12!\n" + "\fcontent_type\x18\x04 \x01(\tR\vcontentType\x12&\n" + "\x0frequest_item_id\x18\x05 \x01(\tR\rrequestItemId\"e\n" + "\vUploadChunk\x127\n" + "\bmetadata\x18\x01 \x01(\v2\x19.credit.v1.UploadMetadataH\x00R\bmetadata\x12\x14\n" + "\x04data\x18\x02 \x01(\fH\x00R\x04dataB\a\n" + "\x05value\"T\n" + "\n" + "DocumentID\x12%\n" + "\x0eapplication_id\x18\x01 \x01(\tR\rapplicationId\x12\x1f\n" + "\vdocument_id\x18\x02 \x01(\tR\n" + "documentId\"c\n" + "\rDownloadChunk\x12\x1b\n" + "\tfile_name\x18\x01 \x01(\tR\bfileName\x12!\n" + "\fcontent_type\x18\x02 \x01(\tR\vcontentType\x12\x12\n" + "\x04data\x18\x03 \x01(\fR\x04data\"<\n" + "\n" + "FieldError\x12\x14\n" + "\x05field\x18\x01 \x01(\tR\x05field\x12\x18\n" + "\amessage\x18\x02 \x01(\tR\amessage\"x\n" + "\n" + "Validation\x12\x14\n" + "\x05valid\x18\x01 \x01(\bR\x05valid\x12%\n" + "\x0emissing_fields\x18\x02 \x03(\tR\rmissingFields\x12-\n" + "\x06errors\x18\x03 \x03(\v2\x15.credit.v1.FieldErrorR\x06errors\"\x95\x02\n" + "\x05Offer\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x1d\n" + "\n" + "product_id\x18\x02 \x01(\tR\tproductId\x12!\n" + "\famount_minor\x18\x03 \x01(\x03R\vamountMinor\x12\x1f\n" + "\vterm_months\x18\x04 \x01(\x05R\n" + "termMonths\x12&\n" + "\x0fannual_rate_bps\x18\x05 \x01(\x05R\rannualRateBps\x122\n" + "\x15monthly_payment_minor\x18\x06 \x01(\x03R\x13monthlyPaymentMinor\x12\x1d\n" + "\n" + "expires_at\x18\a \x01(\tR\texpiresAt\x12\x1e\n" + "\n" + "collateral\x18\b \x01(\tR\n" + "collateral\"\x95\x01\n" + "\bDecision\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" + "\aoutcome\x18\x02 \x01(\tR\aoutcome\x12\x16\n" + "\x06reason\x18\x03 \x01(\tR\x06reason\x12(\n" + "\x06offers\x18\x04 \x03(\v2\x10.credit.v1.OfferR\x06offers\x12\x1d\n" + "\n" + "decided_at\x18\x05 \x01(\tR\tdecidedAt\"\x7f\n" + "\x0fInformationItem\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" + "\x05title\x18\x02 \x01(\tR\x05title\x12\x12\n" + "\x04type\x18\x03 \x01(\tR\x04type\x12\x1a\n" + "\brequired\x18\x04 \x01(\bR\brequired\x12\x16\n" + "\x06status\x18\x05 \x01(\tR\x06status\"\x8f\x01\n" + "\x12InformationRequest\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" + "\amessage\x18\x02 \x01(\tR\amessage\x120\n" + "\x05items\x18\x03 \x03(\v2\x1a.credit.v1.InformationItemR\x05items\x12\x1d\n" + "\n" + "created_at\x18\x04 \x01(\tR\tcreatedAt\"F\n" + "\x0fInformationList\x123\n" + "\x05items\x18\x01 \x03(\v2\x1d.credit.v1.InformationRequestR\x05items\"\xc9\x01\n" + "\x0eManualDecision\x12\x18\n" + "\aoutcome\x18\x01 \x01(\tR\aoutcome\x12\x16\n" + "\x06reason\x18\x02 \x01(\tR\x06reason\x12(\n" + "\x06offers\x18\x03 \x03(\v2\x10.credit.v1.OfferR\x06offers\x12)\n" + "\x10requested_fields\x18\x04 \x03(\tR\x0frequestedFields\x120\n" + "\x05items\x18\x05 \x03(\v2\x1a.credit.v1.InformationItemR\x05items\"T\n" + "\vManualInput\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x125\n" + "\bdecision\x18\x02 \x01(\v2\x19.credit.v1.ManualDecisionR\bdecision\"6\n" + "\fDeclineInput\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x16\n" + "\x06reason\x18\x02 \x01(\tR\x06reason\"8\n" + "\vAcceptInput\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" + "\boffer_id\x18\x02 \x01(\tR\aofferId\"\xb6\x01\n" + "\x04Loan\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12%\n" + "\x0eapplication_id\x18\x02 \x01(\tR\rapplicationId\x12\x19\n" + "\boffer_id\x18\x03 \x01(\tR\aofferId\x12!\n" + "\famount_minor\x18\x04 \x01(\x03R\vamountMinor\x12\x1b\n" + "\tissued_at\x18\x05 \x01(\tR\bissuedAt\x12\x1c\n" + "\tsimulated\x18\x06 \x01(\bR\tsimulated\"\x9e\x01\n" + "\x05Event\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x0e\n" + "\x02at\x18\x02 \x01(\tR\x02at\x12\x1f\n" + "\vfrom_status\x18\x03 \x01(\tR\n" + "fromStatus\x12\x1b\n" + "\tto_status\x18\x04 \x01(\tR\btoStatus\x12\x18\n" + "\amessage\x18\x05 \x01(\tR\amessage\x12\x1d\n" + "\n" + "actor_role\x18\x06 \x01(\tR\tactorRole\"3\n" + "\tEventList\x12&\n" + "\x05items\x18\x01 \x03(\v2\x10.credit.v1.EventR\x05items\"p\n" + "\x05Check\x12\x16\n" + "\x06source\x18\x01 \x01(\tR\x06source\x12\x16\n" + "\x06status\x18\x02 \x01(\tR\x06status\x12\x18\n" + "\asummary\x18\x03 \x01(\tR\asummary\x12\x1d\n" + "\n" + "checked_at\x18\x04 \x01(\tR\tcheckedAt\"\xca\x01\n" + "\rAffordability\x12\x16\n" + "\x06status\x18\x01 \x01(\tR\x06status\x12\x1f\n" + "\vcfads_minor\x18\x02 \x01(\x03R\n" + "cfadsMinor\x12*\n" + "\x11max_payment_minor\x18\x03 \x01(\x03R\x0fmaxPaymentMinor\x12(\n" + "\x10max_amount_minor\x18\x04 \x01(\x03R\x0emaxAmountMinor\x12\x12\n" + "\x04dscr\x18\x05 \x01(\x01R\x04dscr\x12\x16\n" + "\x06reason\x18\x06 \x01(\tR\x06reason\"\x8f\x02\n" + "\aScoring\x12\x14\n" + "\x05score\x18\x01 \x01(\x05R\x05score\x12\x1d\n" + "\n" + "risk_level\x18\x02 \x01(\tR\triskLevel\x12\x18\n" + "\afactors\x18\x03 \x03(\tR\afactors\x12#\n" + "\rmodel_version\x18\x04 \x01(\tR\fmodelVersion\x12(\n" + "\x06checks\x18\x05 \x03(\v2\x10.credit.v1.CheckR\x06checks\x12&\n" + "\x0erecommendation\x18\x06 \x01(\tR\x0erecommendation\x12>\n" + "\raffordability\x18\a \x01(\v2\x18.credit.v1.AffordabilityR\raffordability\"\xb5\x01\n" + "\fRegistryData\x12\x16\n" + "\x06status\x18\x01 \x01(\tR\x06status\x12\x1d\n" + "\n" + "legal_form\x18\x02 \x01(\tR\tlegalForm\x12#\n" + "\rokved_section\x18\x03 \x01(\tR\fokvedSection\x12#\n" + "\rregistered_at\x18\x04 \x01(\tR\fregisteredAt\x12$\n" + "\x0etax_debt_minor\x18\x05 \x01(\x03R\ftaxDebtMinor\"\xa5\x01\n" + "\rCreditHistory\x12\x14\n" + "\x05found\x18\x01 \x01(\bR\x05found\x12*\n" + "\x11overdue_count_12m\x18\x02 \x01(\x05R\x0foverdueCount12m\x12/\n" + "\x14max_overdue_days_24m\x18\x03 \x01(\x05R\x11maxOverdueDays24m\x12!\n" + "\factive_loans\x18\x04 \x01(\x05R\vactiveLoans\"r\n" + "\fCreditReport\x122\n" + "\acompany\x18\x01 \x01(\v2\x18.credit.v1.CreditHistoryR\acompany\x12.\n" + "\x05owner\x18\x02 \x01(\v2\x18.credit.v1.CreditHistoryR\x05owner\"Y\n" + "\vAccountData\x12\x14\n" + "\x05found\x18\x01 \x01(\bR\x05found\x124\n" + "\x16turnover_monthly_minor\x18\x02 \x01(\x03R\x14turnoverMonthlyMinor\"\xe5\x03\n" + "\bFeatures\x12.\n" + "\x13business_age_months\x18\x01 \x01(\x05R\x11businessAgeMonths\x122\n" + "\x15monthly_revenue_minor\x18\x02 \x01(\x03R\x13monthlyRevenueMinor\x124\n" + "\x16monthly_expenses_minor\x18\x03 \x01(\x03R\x14monthlyExpensesMinor\x12=\n" + "\x1bmonthly_debt_payments_minor\x18\x04 \x01(\x03R\x18monthlyDebtPaymentsMinor\x124\n" + "\x16requested_amount_minor\x18\x05 \x01(\x03R\x14requestedAmountMinor\x12\x1f\n" + "\vterm_months\x18\x06 \x01(\x05R\n" + "termMonths\x123\n" + "\bregistry\x18\b \x01(\v2\x17.credit.v1.RegistryDataR\bregistry\x12<\n" + "\rcredit_report\x18\t \x01(\v2\x17.credit.v1.CreditReportR\fcreditReport\x120\n" + "\aaccount\x18\n" + " \x01(\v2\x16.credit.v1.AccountDataR\aaccountJ\x04\b\a\x10\b\"\xf3\x01\n" + "\fScoringInput\x12,\n" + "\x12scoring_request_id\x18\x01 \x01(\tR\x10scoringRequestId\x12%\n" + "\x0eapplication_id\x18\x02 \x01(\tR\rapplicationId\x12/\n" + "\x13application_version\x18\x03 \x01(\x03R\x12applicationVersion\x12/\n" + "\bfeatures\x18\x04 \x01(\v2\x13.credit.v1.FeaturesR\bfeatures\x12,\n" + "\aproduct\x18\x05 \x01(\v2\x12.credit.v1.ProductR\aproduct\"\xc3\x01\n" + "\rScoringOutput\x12,\n" + "\x12scoring_request_id\x18\x01 \x01(\tR\x10scoringRequestId\x12%\n" + "\x0eapplication_id\x18\x02 \x01(\tR\rapplicationId\x12/\n" + "\x13application_version\x18\x03 \x01(\x03R\x12applicationVersion\x12,\n" + "\ascoring\x18\x04 \x01(\v2\x12.credit.v1.ScoringR\ascoring\"\xa2\x01\n" + "\bFeedItem\x12\x0e\n" + "\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" + "\x05title\x18\x02 \x01(\tR\x05title\x12\x12\n" + "\x04body\x18\x03 \x01(\tR\x04body\x12\x16\n" + "\x06reason\x18\x04 \x01(\tR\x06reason\x12\x1d\n" + "\n" + "created_at\x18\x05 \x01(\tR\tcreatedAt\x12%\n" + "\x0eapplication_id\x18\x06 \x01(\tR\rapplicationId\"5\n" + "\bFeedList\x12)\n" + "\x05items\x18\x01 \x03(\v2\x13.credit.v1.FeedItemR\x05items\"\xa6\x03\n" + "\bOverview\x128\n" + "\vapplication\x18\x01 \x01(\v2\x16.credit.v1.ApplicationR\vapplication\x12=\n" + "\x10company_snapshot\x18\x02 \x01(\v2\x12.credit.v1.CompanyR\x0fcompanySnapshot\x12=\n" + "\x10product_snapshot\x18\x03 \x01(\v2\x12.credit.v1.ProductR\x0fproductSnapshot\x12,\n" + "\ascoring\x18\x04 \x01(\v2\x12.credit.v1.ScoringR\ascoring\x12/\n" + "\bdecision\x18\x05 \x01(\v2\x13.credit.v1.DecisionR\bdecision\x121\n" + "\tdocuments\x18\x06 \x03(\v2\x13.credit.v1.DocumentR\tdocuments\x12P\n" + "\x14information_requests\x18\a \x03(\v2\x1d.credit.v1.InformationRequestR\x13informationRequests2\xa6\x04\n" + "\x0eCompanyService\x12<\n" + "\rCreateCompany\x12\x17.credit.v1.CompanyInput\x1a\x12.credit.v1.Company\x12/\n" + "\n" + "GetCompany\x12\r.credit.v1.ID\x1a\x12.credit.v1.Company\x129\n" + "\rListCompanies\x12\x10.credit.v1.Empty\x1a\x16.credit.v1.CompanyList\x12=\n" + "\rUpdateCompany\x12\x18.credit.v1.CompanyUpdate\x1a\x12.credit.v1.Company\x122\n" + "\rVerifyCompany\x12\r.credit.v1.ID\x1a\x12.credit.v1.Company\x128\n" + "\fListProducts\x12\x10.credit.v1.Empty\x1a\x16.credit.v1.ProductList\x12/\n" + "\n" + "GetProduct\x12\r.credit.v1.ID\x1a\x12.credit.v1.Product\x12J\n" + "\tRecommend\x12\x1e.credit.v1.RecommendationInput\x1a\x1d.credit.v1.RecommendationList\x12@\n" + "\tCalculate\x12\x1b.credit.v1.CalculationInput\x1a\x16.credit.v1.Calculation2\x81\n" + "\n" + "\x12ApplicationService\x12B\n" + "\x11CreateApplication\x12\x15.credit.v1.DraftInput\x1a\x16.credit.v1.Application\x127\n" + "\x0eGetApplication\x12\r.credit.v1.ID\x1a\x16.credit.v1.Application\x12D\n" + "\x10ListApplications\x12\x14.credit.v1.PageInput\x1a\x1a.credit.v1.ApplicationPage\x12D\n" + "\x10BankApplications\x12\x14.credit.v1.PageInput\x1a\x1a.credit.v1.ApplicationPage\x12A\n" + "\x10PatchApplication\x12\x15.credit.v1.PatchInput\x1a\x16.credit.v1.Application\x120\n" + "\bValidate\x12\r.credit.v1.ID\x1a\x15.credit.v1.Validation\x12/\n" + "\x06Submit\x12\r.credit.v1.ID\x1a\x16.credit.v1.Application\x12/\n" + "\x06Cancel\x12\r.credit.v1.ID\x1a\x16.credit.v1.Application\x121\n" + "\vGetDecision\x12\r.credit.v1.ID\x1a\x13.credit.v1.Decision\x128\n" + "\x06Accept\x12\x16.credit.v1.AcceptInput\x1a\x16.credit.v1.Application\x12:\n" + "\aDecline\x12\x17.credit.v1.DeclineInput\x1a\x16.credit.v1.Application\x12'\n" + "\x05Issue\x12\r.credit.v1.ID\x1a\x0f.credit.v1.Loan\x12.\n" + "\aHistory\x12\r.credit.v1.ID\x1a\x14.credit.v1.EventList\x12/\n" + "\n" + "GetScoring\x12\r.credit.v1.ID\x1a\x12.credit.v1.Scoring\x128\n" + "\x06Decide\x12\x16.credit.v1.ManualInput\x1a\x16.credit.v1.Application\x121\n" + "\vGetOverview\x12\r.credit.v1.ID\x1a\x13.credit.v1.Overview\x12@\n" + "\x13InformationRequests\x12\r.credit.v1.ID\x1a\x1a.credit.v1.InformationList\x12-\n" + "\x04Feed\x12\x10.credit.v1.Empty\x1a\x13.credit.v1.FeedList\x127\n" + "\rListDocuments\x12\r.credit.v1.ID\x1a\x17.credit.v1.DocumentList\x12?\n" + "\x0eUploadDocument\x12\x16.credit.v1.UploadChunk\x1a\x13.credit.v1.Document(\x01\x129\n" + "\x0eDeleteDocument\x12\x15.credit.v1.DocumentID\x1a\x10.credit.v1.Empty\x12E\n" + "\x10DownloadDocument\x12\x15.credit.v1.DocumentID\x1a\x18.credit.v1.DownloadChunk0\x012O\n" + "\x0eScoringService\x12=\n" + "\bEvaluate\x12\x17.credit.v1.ScoringInput\x1a\x18.credit.v1.ScoringOutputB%Z#robust/contracts/credit/v1;creditv1b\x06proto3"

var (
	file_credit_v1_credit_proto_rawDescOnce sync.Once
	file_credit_v1_credit_proto_rawDescData []byte
)

func file_credit_v1_credit_proto_rawDescGZIP() []byte {
	file_credit_v1_credit_proto_rawDescOnce.Do(func() {
		file_credit_v1_credit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_credit_v1_credit_proto_rawDesc), len(file_credit_v1_credit_proto_rawDesc)))
	})
	return file_credit_v1_credit_proto_rawDescData
}

var file_credit_v1_credit_proto_msgTypes = make([]protoimpl.MessageInfo, 56)
var file_credit_v1_credit_proto_goTypes = []any{(*Empty)(nil), (*ID)(nil), (*CompanyInput)(nil), (*Company)(nil), (*CompanyList)(nil), (*CompanyUpdate)(nil), (*Product)(nil), (*ProductList)(nil), (*Need)(nil), (*RecommendationInput)(nil), (*Recommendation)(nil), (*RecommendationList)(nil), (*CalculationInput)(nil), (*Payment)(nil), (*Calculation)(nil), (*Consent)(nil), (*DraftInput)(nil), (*DraftPatch)(nil), (*Application)(nil), (*PageInput)(nil), (*ApplicationPage)(nil), (*PatchInput)(nil), (*Document)(nil), (*DocumentList)(nil), (*UploadMetadata)(nil), (*UploadChunk)(nil), (*DocumentID)(nil), (*DownloadChunk)(nil), (*FieldError)(nil), (*Validation)(nil), (*Offer)(nil), (*Decision)(nil), (*InformationItem)(nil), (*InformationRequest)(nil), (*InformationList)(nil), (*ManualDecision)(nil), (*ManualInput)(nil), (*DeclineInput)(nil), (*AcceptInput)(nil), (*Loan)(nil), (*Event)(nil), (*EventList)(nil), (*Check)(nil), (*Affordability)(nil), (*Scoring)(nil), (*RegistryData)(nil), (*CreditHistory)(nil), (*CreditReport)(nil), (*AccountData)(nil), (*Features)(nil), (*ScoringInput)(nil), (*ScoringOutput)(nil), (*FeedItem)(nil), (*FeedList)(nil), (*Overview)(nil), nil}
var file_credit_v1_credit_proto_depIdxs = []int32{2, 3, 2, 6, 6, 10, 13, 8, 15, 8, 15, 8, 15, 18, 55, 17, 22, 24, 28, 30, 32, 33, 30, 32, 35, 40, 42, 43, 46, 46, 45, 47, 48, 49, 6, 44, 52, 18, 3, 6, 44, 31, 22, 33, 2, 1, 0, 5, 1, 0, 1, 9, 12, 16, 1, 19, 19, 21, 1, 1, 1, 1, 38, 37, 1, 1, 1, 36, 1, 1, 0, 1, 25, 26, 26, 50, 3, 3, 4, 3, 3, 7, 6, 11, 14, 18, 18, 20, 20, 18, 29, 18, 18, 31, 18, 18, 39, 41, 44, 18, 54, 34, 53, 23, 22, 0, 27, 51, 76, 44, 44, 44, 0}

func init() {
	file_credit_v1_credit_proto_init()
}
func file_credit_v1_credit_proto_init() {
	if File_credit_v1_credit_proto != nil {
		return
	}
	file_credit_v1_credit_proto_msgTypes[9].OneofWrappers = []any{}
	file_credit_v1_credit_proto_msgTypes[17].OneofWrappers = []any{}
	file_credit_v1_credit_proto_msgTypes[25].OneofWrappers = []any{(*UploadChunk_Metadata)(nil), (*UploadChunk_Data)(nil)}
	type x struct{}
	out := protoimpl.TypeBuilder{File: protoimpl.DescBuilder{GoPackagePath: reflect.TypeOf(x{}).PkgPath(), RawDescriptor: unsafe.Slice(unsafe.StringData(file_credit_v1_credit_proto_rawDesc), len(file_credit_v1_credit_proto_rawDesc)), NumEnums: 0, NumMessages: 56, NumExtensions: 0, NumServices: 3}, GoTypes: file_credit_v1_credit_proto_goTypes, DependencyIndexes: file_credit_v1_credit_proto_depIdxs, MessageInfos: file_credit_v1_credit_proto_msgTypes}.Build()
	File_credit_v1_credit_proto = out.File
	file_credit_v1_credit_proto_goTypes = nil
	file_credit_v1_credit_proto_depIdxs = nil
}
