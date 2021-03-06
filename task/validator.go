package task

import (
	"fmt"
)

var TaskValidatorRegistry map[string]func(map[string]interface{}) TaskValidator

func init() {
	if TaskValidatorRegistry == nil {
		TaskValidatorRegistry = make(map[string]func(map[string]interface{}) TaskValidator)
	}
	TaskValidatorRegistry["dummy"] = func(map[string]interface{}) TaskValidator {
		return &DummyTaskValidator{}
	}
}

type TaskValidator interface {
	ValidateTask(t *Task) []error
}

type DummyTaskValidator struct{}

func (d *DummyTaskValidator) ValidateTask(t *Task) []error {
	return nil
}

type Validator struct {
	Validators   map[string]map[string]interface{} `yaml:"validators"`
	myValidators map[string]TaskValidator
}

func (val *Validator) Register() error {
	val.myValidators = make(map[string]TaskValidator)
	for k, v := range val.Validators {
		validator, ok := TaskValidatorRegistry[k]
		if !ok {
			return fmt.Errorf("No config validator for %s", k)
		}
		val.myValidators[k] = validator(v)
	}
	return nil
}

func (v *Validator) ValidateTask(t *Task) []error {
	errs := make([]error, 0)
	err := t.Action.Validate()
	if err != nil {
		errs = append(errs, err)
	}
	validator, ok := v.myValidators[t.Action.RegisteredName()]
	if !ok {
		errs = append(errs, fmt.Errorf("No config for %s", t.Action.RegisteredName()))
		return errs
	}

	for _, err := range validator.ValidateTask(t) {
		errs = append(errs, err)
	}
	return errs
}
