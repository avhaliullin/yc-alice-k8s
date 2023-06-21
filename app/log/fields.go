// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package log

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func FieldJSON(key string, value any) zap.Field {
	return zap.Reflect(key, &jsonObjMarshaller{obj: value})
}

type jsonObjMarshaller struct {
	obj any
}

func (j *jsonObjMarshaller) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(j.obj)
	if err != nil {
		return nil, fmt.Errorf("json marshaling failed: %w", err)
	}
	return bytes, nil
}
