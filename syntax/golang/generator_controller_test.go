package syntax_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
	"github.com/phogolabs/stride/plugin/syntax"
)

var _ = Describe("GeneratorController", func() {
	var generator *syntax.ControllerGenerator

	BeforeEach(func() {
		generator = &syntax.ControllerGenerator{
			Path: tmpdir(),
		}
	})

	Context("when the mode is ControllerGeneratorModeSchema", func() {
		BeforeEach(func() {
			generator.Mode = syntax.ControllerGeneratorModeSchema
		})

		Describe("input", func() {
			It("generates the path parameters", func() {
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user",
							Parameters: codegen.ParameterDescriptorCollection{
								&codegen.ParameterDescriptor{
									Name: "user-id",
									In:   "path",
									ParameterType: &codegen.TypeDescriptor{
										Name:        "string",
										IsPrimitive: true,
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
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Parameters: codegen.ParameterDescriptorCollection{
								&codegen.ParameterDescriptor{
									Name: "name",
									In:   "query",
									ParameterType: &codegen.TypeDescriptor{
										Name:        "string",
										IsPrimitive: true,
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
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Parameters: codegen.ParameterDescriptorCollection{
								&codegen.ParameterDescriptor{
									Name: "X-Partner-ID",
									In:   "header",
									ParameterType: &codegen.TypeDescriptor{
										Name:        "string",
										IsPrimitive: true,
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
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users",
							Name:   "get-user",
							Parameters: codegen.ParameterDescriptorCollection{
								&codegen.ParameterDescriptor{
									Name: "token",
									In:   "cookie",
									ParameterType: &codegen.TypeDescriptor{
										Name:        "string",
										IsPrimitive: true,
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
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "POST",
							Path:   "/users/search",
							Name:   "search-user",
							Requests: codegen.RequestDescriptorCollection{
								&codegen.RequestDescriptor{
									RequestType: &codegen.TypeDescriptor{
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
				Expect(buffer.String()).To(ContainSubstring("Body *SearchQuery `body:\"~\"`"))
			})
		})

		Describe("output", func() {
			It("generates the header parameters", func() {
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user-by-id",
							Responses: codegen.ResponseDescriptorCollection{
								&codegen.ResponseDescriptor{
									Parameters: codegen.ParameterDescriptorCollection{
										&codegen.ParameterDescriptor{
											Name: "X-Partner-ID",
											In:   "header",
											ParameterType: &codegen.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
									ResponseType: &codegen.TypeDescriptor{
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
				generator.Controller = &codegen.ControllerDescriptor{
					Name: "User",
					Operations: codegen.OperationDescriptorCollection{
						&codegen.OperationDescriptor{
							Method: "GET",
							Path:   "/users/{user-id}",
							Name:   "get-user-by-id",
							Responses: codegen.ResponseDescriptorCollection{
								&codegen.ResponseDescriptor{
									Parameters: codegen.ParameterDescriptorCollection{
										&codegen.ParameterDescriptor{
											Name: "X-Partner-ID",
											In:   "header",
											ParameterType: &codegen.TypeDescriptor{
												Name:        "string",
												IsPrimitive: true,
											},
										},
									},
									ResponseType: &codegen.TypeDescriptor{
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
			generator.Mode = syntax.ControllerGeneratorModeAPI
		})
	})

	Context("when the mode is ControllerGeneratorModeSpec", func() {
		BeforeEach(func() {
			generator.Mode = syntax.ControllerGeneratorModeSpec
		})

		//TODO: implement when you introduce spec generation
	})
})
