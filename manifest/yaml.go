package manifest

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type DefaultsSetter interface {
	SetDefaults() error
}

type NameGetter interface {
	GetName() string
}

type NameSetter interface {
	SetName(name string) error
}

func (v Balancers) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (v *Balancers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Balancer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, &v.Endpoints)
}

func (v *Balancer) SetName(name string) error {
	v.Name = name
	return nil
}

func (v BalancerEndpoints) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

// func (v *BalancerEndpoints) UnmarshalYAML(unmarshal func(interface{}) error) error {
//   return unmarshalMapSlice(unmarshal, v)
// }

func (v *BalancerEndpoint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case string:
		v.Target = t
	default:
		return fmt.Errorf("unknown type for endpoint: %T", t)
	}

	return nil
}

func (v *BalancerEndpoint) SetName(name string) error {
	parts := strings.Split(name, "/")

	switch len(parts) {
	case 1:
		v.Port = parts[0]
	case 2:
		v.Port = parts[0]
		v.Protocol = parts[1]
	case 3:
		v.Port = parts[0]
		v.Protocol = parts[1]

		switch parts[2] {
		case "301":
			v.Redirect = v.Target
			v.Target = ""
		default:
			return fmt.Errorf("unknown code: %s", parts[2])
		}
	}

	return nil
}

func (v *Keys) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Key) SetName(name string) error {
	v.Name = name
	return nil
}

func (v Queues) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (v *Queues) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Queue) SetName(name string) error {
	v.Name = name
	return nil
}

func (v Resources) MarshalYAML() (interface{}, error) {
	return marshalMapSlice(v)
}

func (v *Resources) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Resource) SetName(name string) error {
	v.Name = name
	return nil
}

func (v Services) MarshalYAML() (interface{}, error) {
	return marshalMapSlice(v)
}

func (v *Services) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Service) SetName(name string) error {
	v.Name = name
	return nil
}

func (v *ServiceBuild) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case map[interface{}]interface{}:
		type serviceBuild ServiceBuild
		var r serviceBuild
		if err := remarshal(w, &r); err != nil {
			return err
		}
		v.Args = r.Args
		v.Path = r.Path
	case string:
		v.Path = t
	default:
		return fmt.Errorf("unknown type for service build: %T", t)
	}

	return nil
}

func (v *ServiceCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case map[interface{}]interface{}:
		if c, ok := t["development"].(string); ok {
			v.Development = c
		}
		if c, ok := t["test"].(string); ok {
			v.Test = c
		}
		if c, ok := t["production"].(string); ok {
			v.Production = c
		}
	case string:
		v.Development = t
		v.Production = t
	default:
		return fmt.Errorf("unknown type for service command: %T", t)
	}

	return nil
}

func (v *ServiceEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case []interface{}:
		for _, s := range t {
			switch st := s.(type) {
			case []interface{}:
				for _, stv := range st {
					if sv, ok := stv.(string); ok {
						*v = append(*v, sv)
					}
				}
			case string:
				*v = append(*v, st)
			}
		}
	default:
		return fmt.Errorf("unknown type for service environment: %T", t)
	}

	return nil
}

func (v *ServiceHealth) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case map[interface{}]interface{}:
		if w, ok := t["path"].(string); ok {
			v.Path = w
		}
		if w, ok := t["interval"].(int); ok {
			v.Interval = w
		}
		if w, ok := t["timeout"].(int); ok {
			v.Timeout = w
		}
	case string:
		v.Path = t
	default:
		return fmt.Errorf("unknown type for service health: %T", t)
	}

	return nil
}

func (v *ServicePort) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case map[interface{}]interface{}:
		if port := t["port"]; port != nil {
			switch port.(type) {
			case int:
				v.Port = port.(int)
			case string:
				ports, err := strconv.Atoi(port.(string))
				if err != nil {
					return err
				}
				v.Port = ports
			default:
				return fmt.Errorf("invalid port: %v", w)
			}
		}
		if scheme := t["scheme"]; (scheme == nil) || (scheme.(string) == "") {
			v.Scheme = "http"
		} else {
			v.Scheme = scheme.(string)
		}
	case string:
		parts := strings.Split(t, ":")

		switch len(parts) {
		case 1:
			p, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}

			v.Scheme = "http"
			v.Port = p
		case 2:
			p, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			v.Scheme = parts[0]
			v.Port = p
		default:
			return fmt.Errorf("invalid port: %s", t)
		}
	case int:
		v.Scheme = "http"
		v.Port = t
	default:
		return fmt.Errorf("invalid port: %s", t)
	}

	return nil
}

func (v *ServiceScale) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case int:
		v.Count = &ServiceScaleCount{Min: t, Max: t}
	case string:
		var c ServiceScaleCount
		if err := remarshal(w, &c); err != nil {
			return err
		}
		v.Count = &c
	case map[interface{}]interface{}:
		if w, ok := t["count"].(interface{}); ok {
			var c ServiceScaleCount
			if err := remarshal(w, &c); err != nil {
				return err
			}
			v.Count = &c
		}
		if w, ok := t["memory"].(int); ok {
			v.Memory = w
		}
	default:
		return fmt.Errorf("unknown type for service scale: %T", t)
	}

	return nil
}

func (v *ServiceScaleCount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	switch t := w.(type) {
	case int:
		v.Min = t
		v.Max = t
	case string:
		parts := strings.Split(t, "-")

		switch len(parts) {
		case 1:
			i, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}

			v.Min = i

			if !strings.HasSuffix(parts[0], "+") {
				v.Max = i
			}
		case 2:
			i, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}

			j, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			v.Min = i
			v.Max = j
		default:
			return fmt.Errorf("invalid scale: %v", w)
		}
	case map[interface{}]interface{}:
		if min := t["min"]; min != nil {
			switch min.(type) {
			case int:
				v.Min = min.(int)
			case string:
				mins, err := strconv.Atoi(min.(string))
				if err != nil {
					return err
				}
				v.Min = mins
			default:
				return fmt.Errorf("invalid scale: %v", w)
			}
		}
		if max := t["max"]; max != nil {
			switch max.(type) {
			case int:
				v.Max = max.(int)
			case string:
				maxs, err := strconv.Atoi(max.(string))
				if err != nil {
					return err
				}
				v.Max = maxs
			default:
				return fmt.Errorf("invalid scale: %v", w)
			}
		}
	default:
		return fmt.Errorf("invalid scale: %v", w)
	}

	return nil
}

func (v Tables) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (v *Tables) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v *Table) SetName(name string) error {
	v.Name = name
	return nil
}

func (v Timers) MarshalYAML() (interface{}, error) {
	return marshalMapSlice(v)
}

func (v *Timers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshalMapSlice(unmarshal, v)
}

func (v Workflows) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (v *Workflows) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w map[string]map[string]interface{}

	if err := unmarshal(&w); err != nil {
		return err
	}

	for wt, triggers := range w {
		for trigger, stepsi := range triggers {
			steps, ok := stepsi.([]interface{})
			if !ok {
				return fmt.Errorf("could not parse workflow step: %s.%s", wt, trigger)
			}

			wf := Workflow{Type: wt, Trigger: trigger}

			for _, step := range steps {
				switch t := step.(type) {
				case map[interface{}]interface{}:
					if len(t) != 1 {
						return fmt.Errorf("could not parse workflow step: %s.%s", wt, trigger)
					}

					for k, v := range t {
						ks, ok := k.(string)
						if !ok {
							return fmt.Errorf("could not parse workflow step: %s.%s.%v", wt, trigger, k)
						}

						vs, ok := v.(string)
						if !ok {
							return fmt.Errorf("could not parse workflow step: %s.%s.%v", wt, trigger, k)
						}

						wf.Steps = append(wf.Steps, WorkflowStep{Type: ks, Target: vs})
					}
				case string:
					wf.Steps = append(wf.Steps, WorkflowStep{Type: t})
				default:
					return fmt.Errorf("could not parse workflow step: %s.%s", wt, trigger)
				}
			}

			*v = append(*v, wf)
		}
	}

	sort.Slice(*v, func(i, j int) bool {
		vi := (*v)[i]
		vj := (*v)[j]

		if vi.Type == vj.Type {
			return vi.Trigger < vj.Trigger
		}

		return vi.Type < vj.Type
	})

	return nil
}

func remarshal(in, out interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

func marshalMapSlice(in interface{}) (interface{}, error) {
	ms := yaml.MapSlice{}

	iv := reflect.ValueOf(in)

	if iv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice")
	}

	for i := 0; i < iv.Len(); i++ {
		ii := iv.Index(i).Interface()

		if iing, ok := ii.(NameGetter); ok {
			ms = append(ms, yaml.MapItem{
				Key:   iing.GetName(),
				Value: ii,
			})
		}
	}

	return ms, nil
}

func unmarshalMapSlice(unmarshal func(interface{}) error, v interface{}) error {
	rv := reflect.ValueOf(v).Elem()
	vit := rv.Type().Elem()

	var ms yaml.MapSlice

	if err := unmarshal(&ms); err != nil {
		return err
	}

	for _, msi := range ms {
		item := reflect.New(vit).Interface()

		if err := remarshal(msi.Value, item); err != nil {
			return err
		}

		if ds, ok := item.(DefaultsSetter); ok {
			if err := ds.SetDefaults(); err != nil {
				return err
			}
		}

		if ns, ok := item.(NameSetter); ok {
			switch t := msi.Key.(type) {
			case int:
				if err := ns.SetName(fmt.Sprintf("%d", t)); err != nil {
					return err
				}
			case string:
				if err := ns.SetName(t); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown key type: %T", t)
			}
		}

		rv.Set(reflect.Append(rv, reflect.ValueOf(item).Elem()))
	}

	return nil
}
