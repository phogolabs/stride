package codegen_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
)

var _ = Describe("Resolver", func() {
	var spec *codegen.SpecDescriptor

	SchemaAt := func(index int) func() *codegen.TypeDescriptor {
		return func() *codegen.TypeDescriptor {
			By(fmt.Sprintf("getting type descriptor at index %d", index))
			item := spec.Types[index]

			By(fmt.Sprintf("returning type descriptor %+v", item))
			return item
		}
	}

	SchemaElementAt := func(index int) func() *codegen.TypeDescriptor {
		return func() *codegen.TypeDescriptor {
			var (
				get  = SchemaAt(index)
				item = get().Element
			)

			By(fmt.Sprintf("returning type descriptor element %+v", item))
			return item
		}
	}

	ItResolvesPrimitiveType := func(name string, get GetTypeDescriptorFunc) {
		message := fmt.Sprintf("resolve the %s primitive type successfully", name)

		It(message, func() {
			typeSpec := get()

			Expect(typeSpec).NotTo(BeNil())
			Expect(typeSpec.Name).To(Equal(name))
			Expect(typeSpec.IsAlias).To(BeFalse())
			Expect(typeSpec.IsArray).To(BeFalse())
			Expect(typeSpec.IsClass).To(BeFalse())
			Expect(typeSpec.IsEnum).To(BeFalse())
			Expect(typeSpec.IsPrimitive).To(BeTrue())
			Expect(typeSpec.Properties).To(BeEmpty())
		})
	}

	ItResolvesAliasType := func(name string, get GetTypeDescriptorFunc) {
		message := fmt.Sprintf("resolve the %s alias type successfully", name)

		It(message, func() {
			typeSpec := get()

			Expect(typeSpec.Name).To(Equal(name))
			Expect(typeSpec.Element).NotTo(BeNil())
			Expect(typeSpec.IsAlias).To(BeTrue(), "IsAlias")
			Expect(typeSpec.IsArray).To(BeFalse(), "IsArray")
			Expect(typeSpec.IsClass).To(BeFalse(), "IsClass")
			Expect(typeSpec.IsEnum).To(BeFalse(), "IsEnum")
			Expect(typeSpec.IsPrimitive).To(BeFalse(), "IsPrimitive")
			Expect(typeSpec.Properties).To(BeEmpty(), "Properties")
		})
	}

	ItResolvesEnumType := func(name string, get GetTypeDescriptorFunc, values []interface{}) {
		message := fmt.Sprintf("resolve the %s enum type successfully", name)

		It(message, func() {
			typeSpec := get()

			Expect(typeSpec).NotTo(BeNil())
			Expect(typeSpec.Name).To(Equal(name))
			Expect(typeSpec.IsAlias).To(BeFalse())
			Expect(typeSpec.IsArray).To(BeFalse())
			Expect(typeSpec.IsClass).To(BeFalse())
			Expect(typeSpec.IsEnum).To(BeTrue())
			Expect(typeSpec.IsPrimitive).To(BeFalse())
			Expect(typeSpec.Properties).To(BeEmpty())
			Expect(typeSpec.Metadata).To(HaveKeyWithValue("values", values))
		})
	}

	ItResolvesArrayType := func(name string, get GetTypeDescriptorFunc) {
		message := fmt.Sprintf("resolve the %s array type successfully", name)

		It(message, func() {
			typeSpec := get()

			Expect(typeSpec).NotTo(BeNil())
			Expect(typeSpec.Name).To(Equal(name))
			Expect(typeSpec.IsAlias).To(BeFalse())
			Expect(typeSpec.IsArray).To(BeTrue())
			Expect(typeSpec.IsClass).To(BeFalse())
			Expect(typeSpec.IsEnum).To(BeFalse())
			Expect(typeSpec.IsPrimitive).To(BeFalse())
			Expect(typeSpec.Properties).To(BeEmpty())
		})
	}

	ItResolvesObjectType := func(name string, get GetTypeDescriptorFunc) {
		message := fmt.Sprintf("resolve the %s object type successfully", name)

		It(message, func() {
			typeSpec := get()

			Expect(typeSpec).NotTo(BeNil())
			Expect(typeSpec.Name).To(Equal(name))
			Expect(typeSpec.IsAlias).To(BeFalse())
			Expect(typeSpec.IsArray).To(BeFalse())
			Expect(typeSpec.IsClass).To(BeTrue())
			Expect(typeSpec.IsEnum).To(BeFalse())
			Expect(typeSpec.IsPrimitive).To(BeFalse())
			Expect(typeSpec.Properties).NotTo(BeEmpty())
		})
	}

	Describe("Components", func() {
		Describe("Schemas", func() {
			Describe("Integer", func() {
				BeforeEach(func() {
					spec = resolve("schemas-integer.yaml")
					Expect(spec.Types).To(HaveLen(4))
				})

				Describe("Int32", func() {
					ItResolvesPrimitiveType("int32", SchemaElementAt(0))

					ItResolvesAliasType("Int32", SchemaAt(0))
					ItResolvesAliasType("Int32", SchemaElementAt(1))
					ItResolvesAliasType("Int32Ref", SchemaAt(1))
				})

				Describe("int64", func() {
					ItResolvesPrimitiveType("int64", SchemaElementAt(2))

					ItResolvesAliasType("Int64", SchemaAt(2))
					ItResolvesAliasType("Int64", SchemaElementAt(3))
					ItResolvesAliasType("Int64Ref", SchemaAt(3))
				})
			})

			Describe("Number", func() {
				BeforeEach(func() {
					spec = resolve("schemas-number.yaml")
					Expect(spec.Types).To(HaveLen(4))
				})

				Describe("double", func() {
					ItResolvesPrimitiveType("double", SchemaElementAt(0))

					ItResolvesAliasType("Double", SchemaAt(0))
					ItResolvesAliasType("Double", SchemaElementAt(1))
					ItResolvesAliasType("DoubleRef", SchemaAt(1))
				})

				Describe("float", func() {
					ItResolvesPrimitiveType("float", SchemaElementAt(2))

					ItResolvesAliasType("Float", SchemaAt(2))
					ItResolvesAliasType("Float", SchemaElementAt(3))
					ItResolvesAliasType("FloatRef", SchemaAt(3))
				})
			})

			Describe("String", func() {
				BeforeEach(func() {
					spec = resolve("schemas-string.yaml")
					Expect(spec.Types).To(HaveLen(12))
				})

				Describe("binary", func() {
					ItResolvesPrimitiveType("binary", SchemaElementAt(0))

					ItResolvesAliasType("Binary", SchemaAt(0))
					ItResolvesAliasType("Binary", SchemaElementAt(1))
					ItResolvesAliasType("BinaryRef", SchemaAt(1))
				})

				Describe("byte", func() {
					ItResolvesPrimitiveType("byte", SchemaElementAt(2))

					ItResolvesAliasType("Byte", SchemaAt(2))
					ItResolvesAliasType("Byte", SchemaElementAt(3))
					ItResolvesAliasType("ByteRef", SchemaAt(3))
				})

				Describe("date", func() {
					ItResolvesPrimitiveType("date", SchemaElementAt(4))

					ItResolvesAliasType("Date", SchemaAt(4))
					ItResolvesAliasType("Date", SchemaElementAt(5))
					ItResolvesAliasType("DateRef", SchemaAt(5))
				})

				Describe("date-time", func() {
					ItResolvesPrimitiveType("date-time", SchemaElementAt(6))

					ItResolvesAliasType("DateTime", SchemaAt(6))
					ItResolvesAliasType("DateTime", SchemaElementAt(7))
					ItResolvesAliasType("DateTimeRef", SchemaAt(7))
				})

				Describe("string", func() {
					ItResolvesPrimitiveType("string", SchemaElementAt(8))

					ItResolvesAliasType("String", SchemaAt(8))
					ItResolvesAliasType("String", SchemaElementAt(9))
					ItResolvesAliasType("StringRef", SchemaAt(9))
				})

				Describe("uuid", func() {
					ItResolvesPrimitiveType("uuid", SchemaElementAt(10))

					ItResolvesAliasType("UUID", SchemaAt(10))
					ItResolvesAliasType("UUID", SchemaElementAt(11))
					ItResolvesAliasType("UUIDRef", SchemaAt(11))
				})
			})

			Describe("Enum", func() {
				values := []interface{}{"pending", "completed"}

				BeforeEach(func() {
					spec = resolve("schemas-enum.yaml")
					Expect(spec.Types).To(HaveLen(2))
				})

				ItResolvesAliasType("AccountStatus", SchemaAt(0))

				ItResolvesEnumType("TransactionStatus", SchemaElementAt(0), values)
				ItResolvesEnumType("TransactionStatus", SchemaAt(1), values)
			})

			Describe("Array", func() {
				BeforeEach(func() {
					spec = resolve("schemas-array.yaml")
					Expect(spec.Types).To(HaveLen(2))
				})

				ItResolvesArrayType("Array", SchemaAt(0))
				ItResolvesPrimitiveType("string", SchemaElementAt(0))

				ItResolvesAliasType("ArrayRef", SchemaAt(1))
				ItResolvesArrayType("Array", SchemaElementAt(1))
			})

			Describe("Object", func() {
				values := []interface{}{"pending", "completed"}

				BeforeEach(func() {
					spec = resolve("schemas-object.yaml")
					Expect(spec.Types).To(HaveLen(5))
				})

				ItResolvesObjectType("Account", SchemaAt(0))
				ItResolvesObjectType("AccountAddress", SchemaAt(1))
				ItResolvesObjectType("AccountAddressLocation", SchemaAt(2))
				ItResolvesAliasType("AccountRef", SchemaAt(3))
				ItResolvesObjectType("Account", SchemaElementAt(3))
				ItResolvesEnumType("AccountStatus", SchemaAt(4), values)

				It("has a nested property types", func() {
					var property *codegen.PropertyDescriptor

					property = spec.Types[0].Properties[0]
					Expect(property.Name).To(Equal("address"))
					Expect(property.PropertyType.Name).To(Equal("AccountAddress"))
					Expect(property.PropertyType.IsClass).To(BeTrue())
					Expect(property.PropertyType.Properties[2].Name).To(Equal("location"))
					Expect(property.PropertyType.Properties[2].PropertyType.Name).To(Equal("AccountAddressLocation"))
					Expect(property.PropertyType.Properties[2].PropertyType.IsClass).To(BeTrue())

					property = spec.Types[0].Properties[1]
					Expect(property.Name).To(Equal("age"))
					Expect(property.PropertyType.Name).To(Equal("int32"))
					Expect(property.PropertyType.IsPrimitive).To(BeTrue())

					property = spec.Types[0].Properties[2]
					Expect(property.Name).To(Equal("first_name"))
					Expect(property.PropertyType.Name).To(Equal("string"))
					Expect(property.PropertyType.IsPrimitive).To(BeTrue())

					property = spec.Types[0].Properties[3]
					Expect(property.Name).To(Equal("id"))
					Expect(property.PropertyType.Name).To(Equal("uuid"))
					Expect(property.PropertyType.IsPrimitive).To(BeTrue())

					property = spec.Types[0].Properties[4]
					Expect(property.Name).To(Equal("last_name"))
					Expect(property.PropertyType.Name).To(Equal("string"))
					Expect(property.PropertyType.IsPrimitive).To(BeTrue())

					property = spec.Types[0].Properties[5]
					Expect(property.Name).To(Equal("status"))
					Expect(property.PropertyType.Name).To(Equal("AccountStatus"))
					Expect(property.PropertyType.IsEnum).To(BeTrue())
					Expect(property.PropertyType.Metadata).To(HaveKey("values"))
				})
			})
		})

		Describe("Parameters", func() {
		})

		Describe("Responses", func() {
		})

		Describe("Requests", func() {
		})
	})
})
