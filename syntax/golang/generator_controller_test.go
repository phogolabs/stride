package golang_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/fake"
	"github.com/phogolabs/stride/syntax/golang"
)

var _ = Describe("GeneratorController", func() {
	var generator *golang.ControllerGenerator

	BeforeEach(func() {
		reporter := &fake.Reporter{}
		reporter.WithReturns(reporter)

		generator = &golang.ControllerGenerator{
			Path:     tmpdir(),
			Reporter: reporter,
		}
	})

	Context("when the mode is ControllerGeneratorModeSchema", func() {
		BeforeEach(func() {
			generator.Mode = golang.ControllerGeneratorModeSchema
		})

		Describe("input", func() {
			It("generates the path parameters", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user",
							Requests: codedom.RequestDescriptorCollection{
								&codedom.RequestDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "user-id",
											In:   "path",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserInput struct"))
				Expect(buffer.String()).To(ContainSubstring("Path *GetUserInputPath `path:\"~\"`"))

				Expect(buffer.String()).To(ContainSubstring("type GetUserInputPath struct"))
				Expect(buffer.String()).To(ContainSubstring("UserID string `path:\"user-id\" validate:\"-\"`"))
			})

			It("generates the query parameters", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Requests: codedom.RequestDescriptorCollection{
								&codedom.RequestDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "name",
											In:   "query",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserInput struct"))
				Expect(buffer.String()).To(ContainSubstring("Query *GetUserInputQuery `query:\"~\"`"))

				Expect(buffer.String()).To(ContainSubstring("type GetUserInputQuery struct"))
				Expect(buffer.String()).To(ContainSubstring("Name string `query:\"name\" validate:\"-\"`"))
			})

			It("generates the header parameters", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Requests: codedom.RequestDescriptorCollection{
								&codedom.RequestDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "X-Partner-ID",
											In:   "header",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserInput struct"))
				Expect(buffer.String()).To(ContainSubstring("Header *GetUserInputHeader `header:\"~\"`"))

				Expect(buffer.String()).To(ContainSubstring("type GetUserInputHeader struct"))
				Expect(buffer.String()).To(ContainSubstring("XPartnerID string `header:\"X-Partner-ID\" validate:\"-\"`"))
			})

			It("generates the cookie parameters", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Requests: codedom.RequestDescriptorCollection{
								&codedom.RequestDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "token",
											In:   "cookie",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserInput struct"))
				Expect(buffer.String()).To(ContainSubstring("Cookie *GetUserInputCookie `cookie:\"~\"`"))

				Expect(buffer.String()).To(ContainSubstring("type GetUserInputCookie struct"))
				Expect(buffer.String()).To(ContainSubstring("Token string `cookie:\"token\" validate:\"-\"`"))
			})

			It("generates the body", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "POST",
							Path:   "/users/search",
							Name:   "search-user",
							Requests: codedom.RequestDescriptorCollection{
								&codedom.RequestDescriptor{
									RequestType: &codedom.TypeDescriptor{
										Name:       "SearchQuery",
										IsClass:    true,
										IsNullable: true,
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type SearchUserInput struct"))
				Expect(buffer.String()).To(ContainSubstring("Body *SearchQuery `body:\"~\" form:\"~\"`"))
			})
		})

		Describe("output", func() {
			It("generates the header parameters", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user-by-id",
							Responses: codedom.ResponseDescriptorCollection{
								&codedom.ResponseDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "X-Partner-ID",
											In:   "header",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
									ResponseType: &codedom.TypeDescriptor{
										Name:       "User",
										IsClass:    true,
										IsNullable: true,
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserByIDOutput struct"))
				Expect(buffer.String()).To(ContainSubstring("Header *GetUserByIDOutputHeader `header:\"~\"`"))

				Expect(buffer.String()).To(ContainSubstring("type GetUserByIDOutputHeader struct"))
				Expect(buffer.String()).To(ContainSubstring("XPartnerID string `header:\"X-Partner-ID\" validate:\"-\"`"))
			})

			It("generates the body", func() {
				generator.Controller = &codedom.ControllerDescriptor{
					Name: "User",
					Operations: codedom.OperationDescriptorCollection{
						&codedom.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user-by-id",
							Responses: codedom.ResponseDescriptorCollection{
								&codedom.ResponseDescriptor{
									Parameters: codedom.ParameterDescriptorCollection{
										&codedom.ParameterDescriptor{
											Name: "X-Partner-ID",
											In:   "header",
											ParameterType: &codedom.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
									ResponseType: &codedom.TypeDescriptor{
										Name:       "User",
										IsClass:    true,
										IsNullable: true,
									},
								},
							},
						},
					},
				}

				file := generator.Generate()
				Expect(file).NotTo(BeNil())

				buffer := &bytes.Buffer{}
				_, err := file.WriteTo(buffer)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(ContainSubstring("type GetUserByIDOutput struct"))
				Expect(buffer.String()).To(ContainSubstring("Body *User `body:\"~\"`"))
			})
		})
	})

	Context("when the mode is ControllerGeneratorModeAPI", func() {
		BeforeEach(func() {
			generator.Mode = golang.ControllerGeneratorModeAPI
		})
	})

	Context("when the mode is ControllerGeneratorModeSpec", func() {
		BeforeEach(func() {
			generator.Mode = golang.ControllerGeneratorModeSpec
		})

		//TODO: implement when you introduce spec generation
	})
})
