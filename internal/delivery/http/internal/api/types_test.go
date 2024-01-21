package api

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rightJSONTestHelper(ctx context.Context, t *testing.T, s string, value validation.Validatable) {
	if err := json.Unmarshal([]byte(s), value); err != nil {
		require.NoError(t, err)
	}
	assert.NoError(t, value.Validate(ctx, validator.Instance()))
}

func wrongJSONTestHelper(ctx context.Context, t *testing.T, s string, value validation.Validatable) {
	if err := json.Unmarshal([]byte(s), value); err != nil {
		require.NoError(t, err)
	}
	assert.Error(t, value.Validate(ctx, validator.Instance()))
}

func TestAddEducationRequest_Validate(t *testing.T) {
	edJSON := `{
		"number": "1030180354933",
		"date_to": "2015-01-01",
		"date_from": "2011-01-01",
		"issued_institution": "ФГБОУ ВО «Астраханский государственный университет им. В. Н. Татищева»",
		"program": "Связи с общественностью"
	  }`

	var ed AddEducationJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, edJSON, &ed)

	edJSON2 := `{
		"number": "1030180354933",
		"date_to": "2015-01-01",
		"date_from": "2011-01-01",
		"program": "Связи с общественностью"
	  }`
	var ed2 AddEducationRequest
	wrongJSONTestHelper(context.TODO(), t, edJSON2, &ed2)
}

func TestAddContractRequest_Validate(t *testing.T) {
	contractJSON := `{
		"date_from": "2018-01-17",
		"date_to": "2020-01-17",
		"type": "temporary",
		"number": "145678"
	  }`

	var c AddContractJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, contractJSON, &c)
}

func TestAddUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
	}{
		{
			name: "positive #1",
			jsonString: `{
				"first_name": "Alexander",
				"last_name": "Pushkin",
				"middle_name": "Sergeyevich",
				"gender": "male",
				"position": "Novelist",
				"department": "Collegium of Foreign Affairs",
				"place_of_birth": "Moscow",
				"registration_address": "Санкт-Петербург, наб. реки Мойки, 27",
				"residential_address": "Санкт-Петербург, наб. реки Мойки, 27",
				"grade": "1",
				"email": "pushkin@dantes.net",
				"mobile_phone_number": "79999999999",
				"office_phone_number": "123456",
				"insurance": {
					"number": "08336732477"
				},
				"taxpayer": {
					"number": "500100732259"
				},
				"nationality": "russian",
				"foreign_languages": [
				  "english",
				  "german"
				],
				"military": {
				  "rank": "Старший лейтенант",
				  "speciality": "101182",
				  "category": "А2",
				  "comissariat": "Военный комиссариат Петроградского района г. Санкт-Петербурга"
				}
			  }`,
		},
		{
			name: "positive #2",
			jsonString: `{
				"first_name": "Улугбек",
				"last_name": "Акрамов",
				"middle_name": "Рашидович",
				"position": "Строитель",
				"department": "Инженерный отдел",
				"email": "akramovur@rogakopyta.net",
				"mobile_phone_number": "79999999999",
				"office_phone_number": "123456",
				"passports": [
				  {
					"id": 67,
					"number": "AZ0001055",
					"issued_date": "2010-05-23",
					"issued_by": "TOSHKENT SHAHAR IIBB",
					"type": "national",
					"has_scan": true
				  }
				],
				"contracts": [
				  {
					"date_from": "2013-08-01",
					"date_to": "2013-09-04",
					"type": "temporary",
					"number": "12345"
				  }
				],
				"work_permit": {
				  "number": "77121034092",
				  "valid_to": "2013-09-05",
				  "has_scan": true
				},
				"insurance": {
					"number": "08336732477"
				},
				"taxpayer": {
					"number": "500100732259"
				},
				"grade": "2",
				"working_model": "in-office",
				"date_of_birth": "1994-05-26",
				"place_of_birth": "Ташкент",
				"registration_address": "Санкт-Петербург, наб. реки Мойки, 27",
				"residential_address": "Санкт-Петербург, наб. реки Мойки, 27",
				"nationality": "Узбекистан",
				"gender": "male"
			  }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u AddUserJSONRequestBody
			rightJSONTestHelper(context.TODO(), t, tt.jsonString, &u)
		})
	}
}

func TestInsurance_Validate(t *testing.T) {
	insuranceJSON := `{
		"number": "08336732477"
	}`

	var i Insurance
	rightJSONTestHelper(context.TODO(), t, insuranceJSON, &i)
}

func TestLoginJSONRequestBody_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
		wantErr    bool
	}{
		{
			name: "positive",
			jsonString: `{
					"login": "anna@gazneft.ru",
					"password": "pa$$word"
				}`,
			wantErr: false,
		},
		{
			name: "negative #1: bad email format",
			jsonString: `{
					"login": "anna-gazneft.ru",
					"password": "pa$$word"
				}`,
			wantErr: true,
		},
		{
			name: "negative #2: short password",
			jsonString: `{
					"login": "anna@gazneft.ru",
					"password": "pa$"
				}`,
			wantErr: true,
		},
		{
			name: "negative #3: empty password",
			jsonString: `{
					"login": "anna@gazneft.ru",
					"password": ""
				}`,
			wantErr: true,
		},
		{
			name: "negative #4: empty login",
			jsonString: `{
					"login": "",
					"password": "pa$ehdfjguckg"
				}`,
			wantErr: true,
		},
		{
			name: "negative #5: large password",
			jsonString: `{
					"login": "",
					"password": "pa$ehdfjguckppa$"
				}`,
			wantErr: true,
		},
		{
			name: "negative #6: bad email format",
			jsonString: `{
					"login": "A@b@c@example.com",
					"password": "pa$ehdfjguckp$"
				}`,
			wantErr: true,
		},
		{
			name: "negative #7: bad email format",
			jsonString: `{
					"login": "",
					"password": "just\"not\"right@example.com"
				}`,
			wantErr: true,
		},
		{
			name: "negative #8: bad email format",
			jsonString: `{
					"login": "",
					"password": "this is"not\allowed@example.com"
				}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b LoginJSONRequestBody
			errJSON := json.Unmarshal([]byte(tt.jsonString), &b)
			errValidate := b.Validate(context.TODO(), validator.Instance())
			err := errors.Join(errJSON, errValidate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInitChangePasswordJSONRequestBody_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
		wantErr    bool
	}{
		{
			name: "positive",
			jsonString: `{
					"login": "anna@gazneft.ru"
				}`,
			wantErr: false,
		},
		{
			name: "negative #1: bad email format",
			jsonString: `{
					"login": "anna-gazneft.ru"
				}`,
			wantErr: true,
		},
		{
			name: "negative #2: empty login",
			jsonString: `{
					"login": ""
				}`,
			wantErr: true,
		},
		{
			name: "negative #3: bad email format",
			jsonString: `{
					"login": "A@b@c@example.com"
				}`,
			wantErr: true,
		},
		{
			name: "negative #4: bad email format",
			jsonString: `{
					"login": ""
				}`,
			wantErr: true,
		},
		{
			name: "negative #5: bad email format",
			jsonString: `{
					"login": ""
				}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b InitChangePasswordJSONRequestBody
			errJSON := json.Unmarshal([]byte(tt.jsonString), &b)
			errValidate := b.Validate(context.TODO(), validator.Instance())
			err := errors.Join(errJSON, errValidate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	chPswJSON := `{
		"login": "vasyapp@gazneft.ru"
	  }`

	var i InitChangePasswordJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, chPswJSON, &i)
}

func TestChangePasswordRequest_Validate(t *testing.T) {
	chPswJSON := `{
		"key": "0LzQsNC80LAg0LzRi9C70LAg0YDQsNC80YM=",
		"password": "pa$$word"
	  }`

	var c ChangePasswordJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, chPswJSON, &c)
}

func TestAddPassportRequest_Validate(t *testing.T) {
	passportJSON := `{
		"number": "33592222",
		"issued_date": "2016-05-15",
		"issued_by": "Washington D.C. U.S.A.",
		"type": "national"
	  }`

	var p AddPassportJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, passportJSON, &p)
}

func TestTaxpayer_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
	}{
		{
			name: "positive #1",
			jsonString: `{
				"number": "500100732259"
			}`,
		},
		{
			name: "positive #2",
			jsonString: `{
				"number": "1181111110"
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tp Taxpayer
			rightJSONTestHelper(context.TODO(), t, tt.jsonString, &tp)
		})
	}
}

func TestAddTrainingRequest_Validate(t *testing.T) {
	trainingJSON := `{
		"number": "A15/456878-456",
		"issued_institution": "Yandex Practicum",
		"program": "Advanced Go developer",
		"cost": 120000,
		"date_to": "2023-01-10",
		"date_from": "2023-07-10"
	  }`

	var tr AddTrainingJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, trainingJSON, &tr)
}

func TestAddVisaRequest_Validate(t *testing.T) {
	visaJSON := `{
		"number": "33592222",
		"issued_state": "Spain",
		"type": "C1",
		"valid_to": "2017-10-22",
		"valid_from": "2017-09-08",
		"number_entries": "1"
	  }`

	var v AddVisaJSONRequestBody
	rightJSONTestHelper(context.TODO(), t, visaJSON, &v)
}

func TestPatchUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
	}{
		{
			name: "positive #1",
			jsonString: `{
				"position": "Novelist",
				"department": "Collegium of Foreign Affairs",
				"email": "dantes@pushkin.net",
				"mobile_phone_number": "79999999999",
				"foreign_languages": [
				  "english",
				  "german",
				  "french"
				]
			  }`,
		},
		{
			name: "positive #2",
			jsonString: `{
				"position": "Строитель",
				"department": "Инженерный отдел",
				"email": "akramovur@rogakopyta.net",
				"mobile_phone_number": "79999999999",
				"office_phone_number": "123456",
				"contracts": [
				  {
					"date_from": "2013-08-01",
					"date_to": "2013-09-04",
					"type": "temporary",
					"number": "12345"
				  }
				]
			  }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pu PatchUserJSONRequestBody
			rightJSONTestHelper(context.TODO(), t, tt.jsonString, &pu)
		})
	}
}

func TestPatchContractRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
	}{
		{
			name: "positive #1",
			jsonString: `{
				"has_scan": true,
				"date_from": "2018-01-17",
				"type": "temporary",
				"number": "145678"
			  }`,
		},
		{
			name: "positive #2",
			jsonString: `{
				"id": 127,
				"date_to": "2020-01-17",
				"type": "permanent",
				"number": "A/145-S678"
			  }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pc PatchContractJSONRequestBody
			rightJSONTestHelper(context.TODO(), t, tt.jsonString, &pc)
		})
	}
}
