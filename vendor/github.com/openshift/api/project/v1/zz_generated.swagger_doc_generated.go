package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_Project = map[string]string{
	"":         "Projects are the unit of isolation and collaboration in OpenShift. A project has one or more members, a quota on the resources that the project may consume, and the security controls on the resources in the project. Within a project, members may have different roles - project administrators can set membership, editors can create and manage the resources, and viewers can see but not access running containers. In a normal cluster project administrators are not able to alter their quotas - that is restricted to cluster administrators.\n\nListing or watching projects will return only projects the user has the reader role on.\n\nAn OpenShift project is an alternative representation of a Kubernetes namespace. Projects are exposed as editable to end users while namespaces are not. Direct creation of a project is typically restricted to administrators, while end users should use the requestproject resource.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"metadata": "metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "Spec defines the behavior of the Namespace.",
	"status":   "Status describes the current status of a Namespace",
}

func (Project) SwaggerDoc() map[string]string {
	return map_Project
}

var map_ProjectList = map[string]string{
	"":         "ProjectList is a list of Project objects.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"metadata": "metadata is the standard list's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"items":    "Items is the list of projects",
}

func (ProjectList) SwaggerDoc() map[string]string {
	return map_ProjectList
}

var map_ProjectRequest = map[string]string{
	"":            "ProjectRequest is the set of options necessary to fully qualify a project request\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"metadata":    "metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"displayName": "DisplayName is the display name to apply to a project",
	"description": "Description is the description to apply to a project",
}

func (ProjectRequest) SwaggerDoc() map[string]string {
	return map_ProjectRequest
}

var map_ProjectSpec = map[string]string{
	"":           "ProjectSpec describes the attributes on a Project",
	"finalizers": "Finalizers is an opaque list of values that must be empty to permanently remove object from storage",
}

func (ProjectSpec) SwaggerDoc() map[string]string {
	return map_ProjectSpec
}

var map_ProjectStatus = map[string]string{
	"":           "ProjectStatus is information about the current status of a Project",
	"phase":      "Phase is the current lifecycle phase of the project",
	"conditions": "Represents the latest available observations of the project current state.",
}

func (ProjectStatus) SwaggerDoc() map[string]string {
	return map_ProjectStatus
}

// AUTO-GENERATED FUNCTIONS END HERE
