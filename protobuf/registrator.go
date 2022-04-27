package protobuf

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/desc/protoprint"
	"go-kafka-protobuf/srclient"
)

type SchemaRegistrator struct {
	srclient srclient.Client
	printer  *protoprint.Printer
}

func NewSchemaRegistrator(srclient srclient.Client) *SchemaRegistrator {
	printer := &protoprint.Printer{ForceFullyQualifiedNames: true}

	return &SchemaRegistrator{
		srclient: srclient,
		printer:  printer,
	}
}

func (r *SchemaRegistrator) SetKeyCompatibility(ctx context.Context, topic string, compatibility string) {
	// todo: the logic is duplicate here, refactor this to introduce a common function.
	subject := topic + "-key"

	currComp, _ := r.srclient.GetSchemaCompatibility(ctx, subject)
	fmt.Printf("Current key compatibility: %s\n", currComp)

	if currComp != compatibility {
		changeToComp, err := r.srclient.SetSchemaCompatibility(ctx, subject, compatibility)

		if err != nil {
			panic(err)
		}
		fmt.Printf("Updated to key compatibility: %s\n", changeToComp)
	}
}

func (r *SchemaRegistrator) SetValueCompatibility(ctx context.Context, topic string, compatibility string) {
	// todo: the logic is duplicate here, refactor this to introduce a common function.
	subject := topic + "-value"

	currComp, _ := r.srclient.GetSchemaCompatibility(ctx, subject)
	fmt.Printf("Current value compatibility: %s\n", currComp)

	if currComp != compatibility {
		changeToComp, err := r.srclient.SetSchemaCompatibility(ctx, subject, compatibility)

		if err != nil {
			panic(err)
		}
		fmt.Printf("Updated to value compatibility: %s\n", changeToComp)
	}
}

func (r *SchemaRegistrator) RegisterKey(ctx context.Context, topic string, msg interface{}) (int, error) {
	return r.register(ctx, topic+"-key", msg)
}

func (r *SchemaRegistrator) RegisterValue(ctx context.Context, topic string, msg interface{}) (int, error) {
	return r.register(ctx, topic+"-value", msg)
}

func (r *SchemaRegistrator) register(ctx context.Context, topic string, msg interface{}) (int, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return 0, fmt.Errorf("record type must be of proto.Message")
	}

	msgDesc, err := desc.LoadMessageDescriptorForMessage(protoMsg)
	if err != nil {
		return 0, fmt.Errorf("error loading message desciprot for message %w", err)
	}

	fileDesc := msgDesc.GetFile()
	deps := collectFileDescDeps(fileDesc)

	refs := []srclient.Reference{}
	for _, dep := range deps {
		depSchema, err := fileDescriptorToSchemaString(r.printer, dep)
		if err != nil {
			return 0, err
		}

		name := dep.GetName()
		schema, err := r.srclient.CreateSchema(ctx, &srclient.Schema{
			Subject: name,
			Type:    srclient.ProtobufSchemaType,
			Schema:  depSchema,
		})
		if err != nil {
			return 0, fmt.Errorf("Error creating schema: %w", err)
		}

		refs = append(refs, srclient.Reference{
			Name:    name,
			Subject: name,
			Version: schema.Version,
		})
	}

	protoStr, err := fileDescriptorToSchemaString(r.printer, fileDesc)
	if err != nil {
		return 0, err
	}

	schema, err := r.srclient.CreateSchema(ctx, &srclient.Schema{
		Subject:    topic,
		Type:       srclient.ProtobufSchemaType,
		Schema:     protoStr,
		References: refs,
	})
	if err != nil {
		return 0, fmt.Errorf("Error creating schema: %w", err)
	}

	return schema.ID, nil
}

func (r *SchemaRegistrator) Load(ctx context.Context, schemaID int, name string) ([]*desc.FileDescriptor, error) {
	schemaFiles := map[string]string{}
	fileNames := []string{}

	schema, err := r.srclient.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return nil, err
	}

	schemaFiles[name] = schema.Schema
	fileNames = append(fileNames, name)

	for _, dep := range schema.References {
		schema, err := r.srclient.GetSchemaByVersion(ctx, dep.Subject, dep.Version)
		if err != nil {
			return nil, err
		}

		schemaFiles[dep.Name] = schema.Schema
		fileNames = append(fileNames, dep.Name)
	}

	accessor := protoparse.FileContentsFromMap(schemaFiles)

	parser := protoparse.Parser{Accessor: accessor}
	fileDescriptors, err := parser.ParseFiles(fileNames...)
	if err != nil {
		return nil, err
	}

	return fileDescriptors, nil
}

func fileDescriptorToSchemaString(printer *protoprint.Printer, file *desc.FileDescriptor) (string, error) {
	result, err := printer.PrintProtoToString(file)
	if err != nil {
		return "", fmt.Errorf("error converting proto file descriptor to schema string: %w", err)
	}

	return result, nil
}

// reursively collect file descriptor dependencies
func collectFileDescDeps(file *desc.FileDescriptor) []*desc.FileDescriptor {
	var collectDeps func(file *desc.FileDescriptor) []*desc.FileDescriptor
	collectDeps = func(file *desc.FileDescriptor) (deps []*desc.FileDescriptor) {
		for _, dep := range file.GetDependencies() {
			deps = append(deps, dep)
			deps = append(deps, collectDeps(dep)...)
		}

		return
	}

	return collectDeps(file)
}
