package BmFactory

import (
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmHandler"
	"github.com/alfredyang1986/BmPods/BmMiddleware"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/alfredyang1986/BmPods/BmResource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
)

type BmTable struct{}

var BLACKMIRROR_MODEL_FACTORY = map[string]interface{}{
	"BmKid":            BmModel.Kid{},
	"BmApply":          BmModel.Apply{},
	"BmApplicant":      BmModel.Applicant{},
	"BmCategory":       BmModel.Category{},
	"BmCatenode":       BmModel.Catenode{},
	"BmImage":          BmModel.Image{},
	"BmSessioninfo":    BmModel.Sessioninfo{},
	"BmBrand":          BmModel.Brand{},
	"BmReservableitem": BmModel.Reservableitem{},
	"BmStudent":        BmModel.Student{},
	"BmGuardian":       BmModel.Guardian{},
	"BmTeacher":        BmModel.Teacher{},
	"BmRoom":           BmModel.Room{},
	"BmYard":           BmModel.Yard{},
	"BmUnit":           BmModel.Unit{},
	"BmClass":          BmModel.Class{},
	"BmDuty":           BmModel.Duty{},
	"BmTransaction":    BmModel.Transaction{},
	"BmAttachable":     BmModel.Attachable{},
}

var BLACKMIRROR_RESOURCE_FACTORY = map[string]interface{}{
	"BmKidResource":            BmResource.BmKidResource{},
	"BmApplyResource":          BmResource.BmApplyResource{},
	"BmApplicantResource":      BmResource.BmApplicantResource{},
	"BmCategoryResource":       BmResource.BmCategoryResource{},
	"BmCatenodeResource":       BmResource.BmCatenodeResource{},
	"BmImageResource":          BmResource.BmImageResource{},
	"BmSessioninfoResource":    BmResource.BmSessioninfoResource{},
	"BmReservableitemResource": BmResource.BmReservableitemResource{},
	"BmBrandResource":          BmResource.BmBrandResource{},
	"BmStudentResource":        BmResource.BmStudentResource{},
	"BmGuardianResource":       BmResource.BmGuardianResource{},
	"BmTeacherResource":        BmResource.BmTeacherResource{},
	"BmRoomResource":           BmResource.BmRoomResource{},
	"BmYardResource":           BmResource.BmYardResource{},
	"BmUnitResource":           BmResource.BmUnitResource{},
	"BmClassResource":          BmResource.BmClassResource{},
	"BmDutyResource":           BmResource.BmDutyResource{},
	"BmTransactionResource":    BmResource.BmTransactionResource{},
	"BmAttachableResource":     BmResource.BmAttachableResource{},
}

var BLACKMIRROR_STORAGE_FACTORY = map[string]interface{}{
	"BmKidStorage":            BmDataStorage.BmKidStorage{},
	"BmApplyStorage":          BmDataStorage.BmApplyStorage{},
	"BmApplicantStorage":      BmDataStorage.BmApplicantStorage{},
	"BmCategoryStorage":       BmDataStorage.BmCategoryStorage{},
	"BmCatenodeStorage":       BmDataStorage.BmCatenodeStorage{},
	"BmImageStorage":          BmDataStorage.BmImageStorage{},
	"BmSessioninfoStorage":    BmDataStorage.BmSessioninfoStorage{},
	"BmReservableitemStorage": BmDataStorage.BmReservableitemStorage{},
	"BmBrandStorage":          BmDataStorage.BmBrandStorage{},
	"BmStudentStorage":        BmDataStorage.BmStudentStorage{},
	"BmGuardianStorage":       BmDataStorage.BmGuardianStorage{},
	"BmTeacherStorage":        BmDataStorage.BmTeacherStorage{},
	"BmRoomStorage":           BmDataStorage.BmRoomStorage{},
	"BmYardStorage":           BmDataStorage.BmYardStorage{},
	"BmUnitStorage":           BmDataStorage.BmUnitStorage{},
	"BmClassStorage":          BmDataStorage.BmClassStorage{},
	"BmDutyStorage":           BmDataStorage.BmDutyStorage{},
	"BmTransactionStorage":    BmDataStorage.BmTransactionStorage{},
	"BmAttachableStorage":     BmDataStorage.BmAttachableStorage{},
}

var BLACKMIRROR_DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
	"BmSmsDaemon":     BmSms.BmSms{},
}

var BLACKMIRROR_FUNCTION_FACTORY = map[string]interface{}{
	"BmProvinceHandler":        BmHandler.ProvinceHandler{},
	"BmCityHandler":            BmHandler.CityHandler{},
	"BmDistrictHandler":        BmHandler.DistrictHandler{},
	"BmUploadToOssHandler":     BmHandler.UploadToOssHandler{},
	"BmAccountHandler":         BmHandler.AccountHandler{},
	"BmApplicantHandler":       BmHandler.ApplicantHandler{},
	"BmApplicantUpdateHandler": BmHandler.ApplicantUpdateHandler{},
	"BmWeChatHandler":          BmHandler.WeChatHandler{},
	"BmCommonPanicHandle":      BmHandler.CommonPanicHandle{},
	"BmGenerateSmsHandler":     BmHandler.GenerateSmsHandler{},
	"BmVerifiedSmsHandler":     BmHandler.VerifiedSmsHandler{},
	"BmPotentialStudentHandler":     BmHandler.PotentialStudentHandler{},
}

var BLACKMIRROR_MIDDLEWARE_FACTORY = map[string]interface{}{
	"BmCheckTokenMiddleware": BmMiddleware.CheckTokenMiddleware{},
}

func (t BmTable) GetModelByName(name string) interface{} {
	return BLACKMIRROR_MODEL_FACTORY[name]
}

func (t BmTable) GetResourceByName(name string) interface{} {
	return BLACKMIRROR_RESOURCE_FACTORY[name]
}

func (t BmTable) GetStorageByName(name string) interface{} {
	return BLACKMIRROR_STORAGE_FACTORY[name]
}

func (t BmTable) GetDaemonByName(name string) interface{} {
	return BLACKMIRROR_DAEMON_FACTORY[name]
}

func (t BmTable) GetFunctionByName(name string) interface{} {
	return BLACKMIRROR_FUNCTION_FACTORY[name]
}

func (t BmTable) GetMiddlewareByName(name string) interface{} {
	return BLACKMIRROR_MIDDLEWARE_FACTORY[name]
}
