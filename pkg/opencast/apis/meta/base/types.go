/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base

import (
	"encoding/json"
	"time"
)

type Properties map[string]string

type DateTime time.Time

func (dt DateTime) IsZero() bool {
	return time.Time(dt).IsZero()
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte{'"', '"'}, nil
	}
	return time.Time(dt).MarshalJSON()
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' {
		// Opencast represents null dates as blank string. Skip unmarshal and let dt remain as zero value.
		return nil
	}
	return (*time.Time)(dt).UnmarshalJSON(data)
}

type Int int64

func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' {
		// Opencast represents null ints as blank string. Skip unmarshal and let i remain as zero value.
		return nil
	}
	var i2 int64
	if err := json.Unmarshal(data, &i2); err != nil {
		return err
	}
	*i = Int(i2)
	return nil
}

type Flavor string

const (
	DublinCoreEpisodeFlavor = Flavor("dublincore/episode")
	DublinCoreSeriesFlavor  = Flavor("dublincore/series")
	SecurityEpisodeFlavor   = Flavor("security/xacml+episode")
	SecuritySeriesFlavor    = Flavor("security/xacml+series")
	SMILCuttingFlavor       = Flavor("smil/cutting")
	MPEG7SegmentsFlavor     = Flavor("mpeg-7/segments")
)

type Action string

const (
	ReadAction  = Action("read")
	WriteAction = Action("write")
)
