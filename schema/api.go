package schema

import (
	"github.com/saferwall/pe"
	"time"
)

type System struct {
	Platform		   string `json:"platform,omitempty"`
	Architecture 	   string `json:"system_architecture,omitempty"`
	Type 		 	   string `json:"system_type,omitempty"`
	NumberOfProcessors uint32 `json:"number_of_processors,omitempty"`
}

type ResponseObject struct {
	MinioObjectKey string `json:"minio_object_key,omitempty"`
	Sha256         string `json:"sha256,omitempty"`
	FileType       string `json:"file_type"`
}

type Response struct {
	Sha256      string `json:"sha256,omitempty"`
	FileName    string `json:"filename,omitempty"`
	FileSize    int64  `json:"filesize,omitempty"`
	System 		System `json:"system,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

type RequestMalwareScan struct {
	ObjectKey string `json:"minio_object_key,omitempty"`
	Sha256    string `json:"sha256,omitempty"`
	FileType  string `json:"file_type"`
}

type Result struct {
	Status string `json:"status,omitempty"`
}

type FileResponse struct {
	Sha256      string `json:"sha256,omitempty"`
	FileName    string `json:"filename,omitempty"`
	FileSize    int64  `json:"filesize,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

type Submission struct {
	Date     *time.Time `json:"date,omitempty"`
	Filename string     `json:"filename,omitempty"`
	Country  string     `json:"country,omitempty"`
}

type File struct {
	FileKey          string                 `json:"file_key,omitempty"`
	Md5              string                 `json:"md5,omitempty"`
	Sha1             string                 `json:"sha1,omitempty"`
	Sha256           string                 `json:"sha256,omitempty"`
	Sha512           string                 `json:"sha512,omitempty"`
	Ssdeep           string                 `json:"ssdeep,omitempty"`
	Crc32            string                 `json:"crc32,omitempty"`
	Magic            string                 `json:"magic,omitempty"`
	Size             int64                  `json:"size,omitempty"`
	Exif             map[string]string      `json:"exif,omitempty"`
	Tags             map[string]interface{} `json:"tags,omitempty"`
	Packer           []string               `json:"packer,omitempty"`
	FirstSubmission  *time.Time             `json:"first_submission,omitempty"`
	LastSubmission   *time.Time             `json:"last_submission,omitempty"`
	LastScanned      *time.Time             `json:"last_scanned,omitempty"`
	Submissions      []Submission           `json:"submissions,omitempty"`
	SubmissionsCount int64                  `json:"submissions_count"`
	Strings          []StringStruct         `json:"strings,omitempty"`
	MultiAV          map[string]interface{} `json:"multiav,omitempty"`
	Status           int                    `json:"status,omitempty"`
	PE               *pe.File               `json:"pe,omitempty"`
	Histogram        []int                  `json:"histogram,omitempty"`
	ByteEntropy      []int                  `json:"byte_entropy,omitempty"`
	Type             string                 `json:"type,omitempty"`
	Yara             []string               `json:"yara,omitempty"`
}

type StringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}

type ResponseJob struct {
	MinioObjectKey string `json:"minio_object_key,omitempty"`
	Sha256         string `json:"sha256,omitempty"`
	FileType       string `json:"file_type"`
	JobStartStatus bool   `json:"job_start_status"`
}