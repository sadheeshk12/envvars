package envvars

import (
	"errors"
	"fmt"
	"github.com/flemay/envvars/pkg/errorappender"
)

// Validate ensures the Definition is without any error.
// Consumer should always validate before doing any action with the Definition.
func Validate(d *Definition) error {
	return validateDefinitionAndTagNameList(d)
}

func validateDefinition(d *Definition) error {
	errorAppender := errorappender.NewErrorAppender("\n")
	for i, tag := range d.Tags {
		tagErrorAppender := errorappender.NewErrorAppender("; ")
		tagErrorAppender.AppendError(validateTag(tag))
		tagErrorAppender.AppendError(validateTagNameUniqueness(tag.Name, d.Tags))
		tagErrorAppender.AppendError(validateTagUsage(tag.Name, d.Envvars))
		errorAppender.AppendError(tagErrorAppender.Wrap(fmt.Sprintf("Tag '%s' (#%d): ", tag.Name, i+1)))
	}
	for i, ev := range d.Envvars {
		evErrorAppender := errorappender.NewErrorAppender("; ")
		evErrorAppender.AppendError(validateEnvvar(ev, d.Tags))
		evErrorAppender.AppendError(validateEnvvarNameUniqueness(ev.Name, d.Envvars))
		errorAppender.AppendError(evErrorAppender.Wrap(fmt.Sprintf("Envvar '%s' (#%d): ", ev.Name, i+1)))
	}

	return errorAppender.Error()
}

func validateDefinitionAndTagNameList(d *Definition, tagNames ...string) error {
	errorAppender := errorappender.NewErrorAppender("\n")
	errorAppender.AppendError(validateDefinition(d))
	errorAppender.AppendError(validateTagNameList(tagNames, d.Tags))
	return errorAppender.Error()
}

func validateEnvvar(ev *Envvar, tags TagCollection) error {
	errorAppender := errorappender.NewErrorAppender("; ")
	if ev.Name == "" {
		errorAppender.AppendString("name cannot be blank")
	}
	if ev.Desc == "" {
		errorAppender.AppendString("desc cannot be blank")
	}
	errorAppender.AppendError(validateTagNameList(ev.Tags, tags))

	return errorAppender.Error()
}

func validateEnvvarNameUniqueness(name string, c EnvvarCollection) error {
	if len(c.GetAll(name)) > 1 {
		return errors.New("name must be unique")
	}
	return nil
}

func validateTag(t *Tag) error {
	errorAppender := errorappender.NewErrorAppender("; ")
	if t.Name == "" {
		errorAppender.AppendString("name cannot be blank")
	}
	if t.Desc == "" {
		errorAppender.AppendString("desc cannot be blank")
	}
	return errorAppender.Error()
}

func validateTagNameUniqueness(name string, tags TagCollection) error {
	if len(tags.GetAll(name)) > 1 {
		return errors.New("name must be unique")
	}
	return nil
}

func validateTagUsage(name string, c EnvvarCollection) error {
	if len(c.WithTag(name)) == 0 {
		return errors.New("is not used")
	}
	return nil
}

func validateTagNameList(names []string, tags TagCollection) error {
	errorAppender := errorappender.NewErrorAppender("; ")
	counts := make(map[string]int)
	for _, name := range names {
		if name == "" {
			errorAppender.AppendString("tag '' is empty")
			continue
		}
		if tags.Get(name) == nil {
			errorAppender.AppendError(fmt.Errorf("tag '%v' is not defined", name))
		}
		counts[name] = counts[name] + 1
	}
	for name, counts := range counts {
		if counts > 1 {
			errorAppender.AppendError(fmt.Errorf("tag '%v' is duplicated", name))
		}
	}
	return errorAppender.Error()
}
