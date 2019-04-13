package codegen_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
)

var _ = Describe("Resolver", func() {
	var (
		resolver *codegen.Resolver
		spec     *codegen.SpecDescriptor
	)

	BeforeEach(func() {
		resolver = &codegen.Resolver{}
		Expect(resolver).NotTo(BeNil())
	})

	JustBeforeEach(func() {
		swagger := load("../fixture/swagger.yaml")
		spec = resolver.Resolve(swagger)
	})

	It("resolves the schemas successfully", func() {
		Expect(spec.Schemas).To(HaveLen(10))

		for index, schema := range spec.Schemas {
			switch index {
			case 0:
				Expect(schema.Name).To(Equal("Account"))
				Expect(schema.Properties).To(HaveLen(4))

				Expect(schema.Properties[0].Name).To(Equal("balance"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("#/components/schemas/Balance"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeFalse())

				Expect(schema.Properties[1].Name).To(Equal("bank_id"))
				Expect(schema.Properties[1].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[1].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[1].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[2].Name).To(Equal("id"))
				Expect(schema.Properties[2].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[2].PropertyType.Name).To(Equal("uuid"))
				Expect(schema.Properties[2].PropertyType.IsPrimitive).To(BeTrue())
			case 1:
				Expect(schema.Name).To(Equal("AccountArray"))
				Expect(schema.IsArray).To(BeTrue())
				Expect(schema.Properties).To(HaveLen(4))
			case 2:
				Expect(schema.Name).To(Equal("AccountArrayOutput"))
				Expect(schema.Properties).To(HaveLen(1))

				Expect(schema.Properties[0].Name).To(Equal("accounts"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("#/components/schemas/Account"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeFalse())
			case 3:
				Expect(schema.Name).To(Equal("AccountOutput"))
				Expect(schema.Properties).To(HaveLen(1))

				Expect(schema.Properties[0].Name).To(Equal("account"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("#/components/schemas/Account"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeFalse())
			case 4:
				Expect(schema.Name).To(Equal("Balance"))
				Expect(schema.Properties).To(HaveLen(2))

				Expect(schema.Properties[0].Name).To(Equal("amount"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("double"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[1].Name).To(Equal("currency"))
				Expect(schema.Properties[1].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[1].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[1].PropertyType.IsPrimitive).To(BeTrue())
			case 5:
				Expect(schema.Name).To(Equal("Customer"))
				Expect(schema.Properties).To(HaveLen(5))

				Expect(schema.Properties[0].Name).To(Equal("date_of_birth"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("date-time"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[1].Name).To(Equal("email"))
				Expect(schema.Properties[1].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[1].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[1].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[2].Name).To(Equal("employment_status"))
				Expect(schema.Properties[2].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[2].PropertyType.Name).To(Equal("#/components/schemas/EmploymentStatus"))
				Expect(schema.Properties[2].PropertyType.IsPrimitive).To(BeFalse())

				Expect(schema.Properties[3].Name).To(Equal("first_name"))
				Expect(schema.Properties[3].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[3].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[3].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[4].Name).To(Equal("last_name"))
				Expect(schema.Properties[4].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[4].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[4].PropertyType.IsPrimitive).To(BeTrue())
			case 6:
				Expect(schema.Name).To(Equal("CustomerOutput"))
				Expect(schema.Properties).To(HaveLen(1))

				Expect(schema.Properties[0].Name).To(Equal("customer_id"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("uuid"))
			case 7:
				Expect(schema.Name).To(Equal("EmploymentStatus"))
				Expect(schema.Properties).To(HaveLen(2))
				Expect(schema.Properties[0].Name).To(Equal("employed"))
				Expect(schema.Properties[1].Name).To(Equal("unemployed"))
			case 8:
				Expect(schema.Name).To(Equal("Transaction"))
				Expect(schema.Properties).To(HaveLen(8))

				Expect(schema.Properties[0].Name).To(Equal("balance"))
				Expect(schema.Properties[0].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[0].PropertyType.Name).To(Equal("#/components/schemas/Balance"))
				Expect(schema.Properties[0].PropertyType.IsPrimitive).To(BeFalse())

				Expect(schema.Properties[1].Name).To(Equal("category"))
				Expect(schema.Properties[1].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[1].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[1].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[2].Name).To(Equal("created_at"))
				Expect(schema.Properties[2].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[2].PropertyType.Name).To(Equal("date-time"))
				Expect(schema.Properties[2].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[3].Name).To(Equal("description"))
				Expect(schema.Properties[3].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[3].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[3].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[4].Name).To(Equal("id"))
				Expect(schema.Properties[4].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[4].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[4].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[5].Name).To(Equal("merchant"))
				Expect(schema.Properties[5].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[5].PropertyType.Name).To(Equal("string"))
				Expect(schema.Properties[5].PropertyType.IsPrimitive).To(BeTrue())

				Expect(schema.Properties[6].Name).To(Equal("settled_at"))
				Expect(schema.Properties[6].PropertyType).NotTo(BeNil())
				Expect(schema.Properties[6].PropertyType.Name).To(Equal("date-time"))
			case 9:
				Expect(schema.Name).To(Equal("TransactionArray"))
				Expect(schema.Properties).To(HaveLen(8))
				Expect(schema.IsArray).To(BeTrue())
			}
		}
	})

	It("resolves the requests successfully", func() {
		Expect(spec.RequestBodies).To(HaveLen(1))
		Expect(spec.RequestBodies[0].Name).To(Equal("CustomerInput"))
		Expect(spec.RequestBodies[0].Description).To(Equal("Customer's details"))
		Expect(spec.RequestBodies[0].Contents).To(HaveLen(1))
		Expect(spec.RequestBodies[0].Contents[0].Name).To(Equal("application/json"))
		Expect(spec.RequestBodies[0].Contents[0].ContentType).NotTo(BeNil())
		Expect(spec.RequestBodies[0].Contents[0].ContentType.Name).To(Equal("#/components/schemas/Customer"))
	})

	It("resolves the responses successfully", func() {
		Expect(spec.Responses).To(HaveLen(1))
		Expect(spec.Responses[0].Code).To(Equal(0))
		Expect(spec.Responses[0].Name).To(Equal("CustomerOutput"))
		Expect(spec.Responses[0].Description).To(Equal("Customer's details"))

		Expect(spec.Responses[0].Headers).To(HaveLen(1))
		Expect(spec.Responses[0].Headers[0].Name).To(Equal("X-Rate-Limit"))
		Expect(spec.Responses[0].Headers[0].HeaderType.Name).To(Equal("int32"))
		Expect(spec.Responses[0].Headers[0].HeaderType.IsPrimitive).To(BeTrue())

		Expect(spec.Responses[0].Contents).To(HaveLen(1))
		Expect(spec.Responses[0].Contents[0].Name).To(Equal("application/json"))
		Expect(spec.Responses[0].Contents[0].ContentType).NotTo(BeNil())
		Expect(spec.Responses[0].Contents[0].ContentType.Name).To(BeEmpty())
		Expect(spec.Responses[0].Contents[0].ContentType.IsClass).To(BeTrue())
		Expect(spec.Responses[0].Contents[0].ContentType.Properties).To(HaveLen(1))
		Expect(spec.Responses[0].Contents[0].ContentType.Properties[0].Name).To(Equal("account_id"))
		Expect(spec.Responses[0].Contents[0].ContentType.Properties[0].PropertyType.Name).To(Equal("uuid"))
	})

	It("resolves the parameters successfully", func() {
		Expect(spec.Parameters).To(HaveLen(1))
		Expect(spec.Parameters[0].Name).To(Equal("example"))
		Expect(spec.Parameters[0].Description).To(Equal("Some Parameter"))
		Expect(spec.Parameters[0].ParameterType.Name).To(Equal("int32"))
		Expect(spec.Parameters[0].ParameterType.IsPrimitive).To(BeTrue())
	})

	It("resolves the controllers successfully", func() {
		Expect(spec.Controllers).To(HaveLen(3))
		Expect(spec.Controllers[0].Name).To(Equal("account"))
		Expect(spec.Controllers[0].Operations).To(HaveLen(2))

		Expect(spec.Controllers[0].Operations[0].Name).To(Equal("getAccountById"))
		Expect(spec.Controllers[0].Operations[0].Parameters).To(HaveLen(1))
		Expect(spec.Controllers[0].Operations[0].Parameters[0].Name).To(Equal("accountId"))
		Expect(spec.Controllers[0].Operations[0].RequestBody).To(BeNil())
		Expect(spec.Controllers[0].Operations[0].Responses).To(HaveLen(1))
		Expect(spec.Controllers[0].Operations[0].Responses[0].Code).To(Equal(200))
		Expect(spec.Controllers[0].Operations[0].Responses[0].Contents).To(HaveLen(1))
		Expect(spec.Controllers[0].Operations[0].Responses[0].Contents[0].Name).To(Equal("application/json"))
		Expect(spec.Controllers[0].Operations[0].Responses[0].Contents[0].ContentType).NotTo(BeNil())
		Expect(spec.Controllers[0].Operations[0].Responses[0].Contents[0].ContentType.Name).To(Equal("#/components/schemas/AccountOutput"))

		Expect(spec.Controllers[0].Operations[1].Name).To(Equal("getAccounts"))
		Expect(spec.Controllers[0].Operations[1].Parameters).To(HaveLen(0))
		Expect(spec.Controllers[0].Operations[1].RequestBody).To(BeNil())
		Expect(spec.Controllers[0].Operations[1].Responses).To(HaveLen(1))
		Expect(spec.Controllers[0].Operations[1].Responses[0].Code).To(Equal(200))
		Expect(spec.Controllers[0].Operations[1].Responses[0].Contents).To(HaveLen(1))
		Expect(spec.Controllers[0].Operations[1].Responses[0].Contents[0].Name).To(Equal("application/json"))
		Expect(spec.Controllers[0].Operations[1].Responses[0].Contents[0].ContentType).NotTo(BeNil())
		Expect(spec.Controllers[0].Operations[1].Responses[0].Contents[0].ContentType.Name).To(Equal("#/components/schemas/AccountArrayOutput"))
	})
})
