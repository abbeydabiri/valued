package database

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type Profile struct {
	Crud
}

func (this *Profile) GetName() string {
	return "profile"
}

func (this *Profile) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Profile) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["dob"] = ""
	fields["gender"] = ""
	fields["nationality"] = ""

	fields["firstname"] = ""
	fields["lastname"] = ""
	fields["referrer"] = ""
	fields["image"] = ""

	fields["employercontrol"] = ""

	fields["email"] = ""
	fields["emailsecondary"] = ""
	fields["emailalternate"] = ""
	fields["phone"] = ""
	fields["phonecode"] = ""
	fields["website"] = ""

	fields["categorycontrol"] = ""
	fields["subcategorycontrol"] = ""

	fields["keywords"] = ""

	fields["merchantterms"] = ""
	fields["merchantbenefits"] = ""
	fields["activationdate"] = ""

	fields["username"] = ""
	fields["password"] = ""
	fields["pincode"] = ""

	fields["status"] = ""

	fields["company"] = "Yes"

	return fields
}

func (this *Profile) EncryptPassword(cString string, curdb Database) string {

	byteEncrypted := curdb.Encrypt([]byte(cString))
	return base64.StdEncoding.EncodeToString(byteEncrypted)
	// return base64.StdEncoding.EncodeToString([]byte(cString))
}

func (this *Profile) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["code"] = "main"
	xDoc["title"] = "VALUED PEOPLE"
	xDoc["firstname"] = ""
	xDoc["lastname"] = ""

	xDoc["activationdate"] = ""
	xDoc["workflow"] = "registered"
	xDoc["status"] = "active"
	control := this.Create("root", xDoc, curdb)

	xDoc["categorycontrol"] = control
	xDoc["subcategorycontrol"] = control

	xDoc["employercontrol"] = control

	xDoc["username"] = "root@localhost.com"
	xDoc["password"] = "toor"
	xDoc["pincode"] = "4321"
	xDoc["company"] = ""
	this.Update("root", xDoc, curdb)

}

func (this *Profile) VerifyLogin(sRole, cUsername, cPassword string, curdb Database) (mapRes map[string]interface{}) {

	sFrontendColumn := fmt.Sprintf(` ,'%s' as role`, sRole)
	sFrontend := fmt.Sprintf(`AND profile.control in ( 
						select profilerole.profilecontrol from profilerole 
						left join role on role.control = profilerole.rolecontrol
						where lower(role.code) = lower('%s')
					)`, sRole)

	sql := `select distinct profile.code as code, profile.control as control, profile.title as title, 
					profile.description as description, profile.workflow as workflow, profile.status as status, 
					profile.createdate as createdate, profile.createdby as createdby, 
					profile.updatedate as updatedate, profile.updatedby as updatedby,

					profile.firstname as firstname, profile.lastname as lastname, profile.referrer as referrer, 
					profile.gender as gender, profile.dob as dob, profile.nationality as nationality, 

					profile.website as website, profile.phonecode as phonecode, profile.phone as phone, profile.email as email, 
					profile.emailsecondary as emailsecondary, profile.emailalternate as emailalternate,

					profile.image as image, profile.keywords as keywords, profile.merchantterms as merchantterms, 
					profile.merchantbenefits as merchantbenefits, 

					profile.pincode as pincode, profile.username as username, profile.password as password,
					profile.company as company,

					profile.employercontrol as employercontrol, employer.title as employertitle, employer.code as employercode,
					employer.email as employeremail, employer.createdate as employercreatedate,
					
					profile.categorycontrol as categorycontrol, category.title as categorytitle,
                    profile.subcategorycontrol as subcategorycontrol, subcategory.title as subcategorytitle

                    %s

			from profile
            
            left join category on profile.categorycontrol = category.control
            left join category as subcategory on profile.subcategorycontrol = subcategory.control
                       
            left join profile as employer on profile.employercontrol = employer.control 

			where profile.username = '%s' AND profile.password = '%s' AND profile.status = 'active'  %s
			`
	sql = fmt.Sprintf(sql, sFrontendColumn, cUsername, cPassword, sFrontend)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Profile) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Profile) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Profile) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select distinct profile.code as code, profile.control as control, profile.title as title, 
					profile.description as description, profile.workflow as workflow, profile.status as status, 
					profile.createdate as createdate, profile.createdby as createdby, 
					profile.updatedate as updatedate, profile.updatedby as updatedby,

					profile.firstname as firstname, profile.lastname as lastname, profile.referrer as referrer, 
					profile.gender as gender, profile.dob as dob, profile.nationality as nationality, 

					profile.website as website, profile.phonecode as phonecode, profile.phone as phone, profile.email as email, 
					profile.emailsecondary as emailsecondary, profile.emailalternate as emailalternate,

					profile.image as image, profile.keywords as keywords, profile.merchantterms as merchantterms, 
					profile.merchantbenefits as merchantbenefits, 

					profile.pincode as pincode, profile.username as username, profile.password as password,
					profile.company as company,

					profile.employercontrol as employercontrol, employer.title as employertitle,  employer.code as employercode,
					employer.email as employeremail, employer.createdate as employercreatedate,
					
					profile.categorycontrol as categorycontrol, category.title as categorytitle,
                    profile.subcategorycontrol as subcategorycontrol, subcategory.title as subcategorytitle


			from profile
            

            left join category on profile.categorycontrol = category.control
            left join category as subcategory on profile.subcategorycontrol = subcategory.control
                       
            left join profile as employer on profile.employercontrol = employer.control 
            
            where profile.%s = '%s' 
					
			order by firstname, lastname, title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Profile) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`profile.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherestatus := ""
	if xDocrequest["status"] != nil && xDocrequest["status"].(string) != "" {
		wherestatus = fmt.Sprintf(`profile.status = '%s' AND`, xDocrequest["status"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`lower(profile.title) like lower('%%%s%%') AND`, xDocrequest["title"])
	}

	wherefirstname := ""
	if xDocrequest["firstname"] != nil && xDocrequest["firstname"].(string) != "" {
		wherefirstname = fmt.Sprintf(`lower(profile.firstname) like lower('%%%s%%') AND`, xDocrequest["firstname"])
	}

	wherelastname := ""
	if xDocrequest["lastname"] != nil && xDocrequest["lastname"].(string) != "" {
		wherelastname = fmt.Sprintf(`lower(profile.lastname) like lower('%%%s%%') AND`, xDocrequest["lastname"])
	}

	wherereferrer := ""
	if xDocrequest["referrer"] != nil && xDocrequest["referrer"].(string) != "" {
		wherereferrer = fmt.Sprintf(`lower(profile.referrer) like lower('%%%s%%') AND`, xDocrequest["referrer"])
	}

	whereemail := ""
	if xDocrequest["email"] != nil && xDocrequest["email"].(string) != "" {
		whereemail = fmt.Sprintf(`lower(profile.email) like lower('%%%s%%') AND`, xDocrequest["email"])
	}

	whereusername := ""
	if xDocrequest["username"] != nil && xDocrequest["username"].(string) != "" {
		whereusername = fmt.Sprintf(`lower(profile.username) like lower('%%%s%%') AND`, xDocrequest["username"])
	}

	whereprofilecontrol := ""
	if xDocrequest["profilecontrol"] != nil && len(xDocrequest["profilecontrol"].([]string)) > 0 {
		whereprofilecontrol = strings.Join(xDocrequest["profilecontrol"].([]string), `","`)
		whereprofilecontrol = fmt.Sprintf(`profile.control in ("%s") AND`, whereprofilecontrol)
	}

	whereemployertitle := ""
	if xDocrequest["employertitle"] != nil && xDocrequest["employertitle"].(string) != "" {
		whereemployertitle = fmt.Sprintf(`lower(employer.title) like lower('%%%s%%') AND`, xDocrequest["employertitle"])
	}

	whereemployercode := ""
	if xDocrequest["employercode"] != nil && xDocrequest["employercode"].(string) != "" {
		whereemployercode = fmt.Sprintf(`lower(employer.code) = lower('%s') AND`, xDocrequest["employercode"])
	}

	wherecompany := "lower(profile.company) != lower('Yes') AND"
	if xDocrequest["company"] != nil && xDocrequest["company"].(string) != "" {
		wherecompany = fmt.Sprintf(`lower(profile.company) = lower('Yes') AND`)
	}

	whereemployercontrol := ""
	if xDocrequest["employercontrol"] != nil && xDocrequest["employercontrol"].(string) != "" {
		whereemployercontrol = fmt.Sprintf(`profile.employercontrol = '%s' AND`, xDocrequest["employercontrol"])
	}

	wherecategoryname := ""
	if xDocrequest["categoryname"] != nil && xDocrequest["categoryname"].(string) != "" {
		wherecategoryname = fmt.Sprintf(`category.title like '%%%s%%' AND`, xDocrequest["categoryname"])
	}

	wherecategory := ""
	if xDocrequest["category"] != nil && xDocrequest["category"].(string) != "" {
		wherecategory = fmt.Sprintf(`category.control = '%s' AND`, xDocrequest["category"])
	}

	wheresubcategory := ""
	if xDocrequest["subcategory"] != nil && xDocrequest["subcategory"].(string) != "" {
		wheresubcategory = fmt.Sprintf(`profile.subcategorycontrol = '%s' AND`, xDocrequest["subcategory"])
	}

	wherekeywords := ""
	if xDocrequest["keywords"] != nil && xDocrequest["keywords"].(string) != "" {
		wherekeywords = fmt.Sprintf(`lower(profile.keywords) like lower('%%%s%%') AND`, xDocrequest["keywords"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`lower(profile.code) like lower('%%%s%%') AND`, xDocrequest["code"])
	}

	whererole := ""
	whereroleColumn := ""
	if xDocrequest["role"] != nil && xDocrequest["role"].(string) != "" {
		whererole = `profile.control in ( 
						select profilerole.profilecontrol from profilerole 
						left join role on role.control = profilerole.rolecontrol
						where lower(role.code) = lower('%s')
					)   AND`
		whererole = fmt.Sprintf(whererole, xDocrequest["role"])
		whereroleColumn = fmt.Sprintf(` ,'%s' as role`, xDocrequest["role"])
	}

	sqlFields := `
					distinct 

					profile.code as code, profile.control as control, profile.title as title, 
					profile.description as description, profile.workflow as workflow, profile.status as status, 
					profile.createdate as createdate, profile.createdby as createdby, 
					profile.updatedate as updatedate, profile.updatedby as updatedby,

					profile.firstname as firstname, profile.lastname as lastname, profile.referrer as referrer, 
					profile.gender as gender, profile.dob as dob, profile.nationality as nationality, 

					profile.website as website, profile.phonecode as phonecode, profile.phone as phone, profile.email as email, 
					profile.emailsecondary as emailsecondary, profile.emailalternate as emailalternate,

					profile.image as image, profile.keywords as keywords, profile.merchantterms as merchantterms, 
					profile.merchantbenefits as merchantbenefits, 

					profile.pincode as pincode, profile.username as username, profile.password as password,
					profile.company as company,

					profile.employercontrol as employercontrol, employer.title as employertitle,  employer.code as employercode,
					employer.email as employeremail, employer.createdate as employercreatedate,

					profile.categorycontrol as categorycontrol, category.title as categorytitle,
					profile.subcategorycontrol as subcategorycontrol, subcategory.title as subcategorytitle

				` + whereroleColumn

	sqlOrderByLimit := `order by profile.firstname, profile.lastname, profile.title limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(distinct profile.control) as paginationtotal	`
	}

	sql := `select %s

			from profile
                                               
            left join category on profile.categorycontrol = category.control
            left join category as subcategory on profile.subcategorycontrol = subcategory.control
                       
            left join profile as employer on profile.employercontrol = employer.control 
            
            
            where  %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s 
					
			profile.control != '' ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherestatus, wheretitle, wherecompany, wherefirstname, wherelastname,
		wherereferrer, whereemail, whereusername, whereprofilecontrol, whereemployertitle, whereemployercode, whereemployercontrol,
		wherecategoryname, wherecategory, wheresubcategory, wherekeywords, wherecode, whererole, limit, offset)
	mapRes, _ = curdb.Query(sql)
	// println(sql)
	return mapRes
}

func (this *Profile) SearchEmployee(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	wherefirstname := ""
	if xDocrequest["firstname"] != nil && xDocrequest["firstname"].(string) != "" {
		wherefirstname = fmt.Sprintf(`lower(profile.firstname) like lower('%%%s%%') AND`, xDocrequest["firstname"])
	}

	wherelastname := ""
	if xDocrequest["lastname"] != nil && xDocrequest["lastname"].(string) != "" {
		wherelastname = fmt.Sprintf(`lower(profile.lastname) like lower('%%%s%%') AND`, xDocrequest["lastname"])
	}

	whereemail := ""
	if xDocrequest["email"] != nil && xDocrequest["email"].(string) != "" {
		whereemail = fmt.Sprintf(`lower(profile.email) like lower('%%%s%%') AND`, xDocrequest["email"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`lower(profile.code) like lower('%%%s%%') AND`, xDocrequest["code"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`profile.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherescheme := ""
	if xDocrequest["scheme"] != nil && xDocrequest["scheme"].(string) != "" {
		wherescheme = fmt.Sprintf(`subscription.schemecontrol =  '%s' AND`, xDocrequest["scheme"])
	}

	whereemployer := ""
	if xDocrequest["employer"] != nil && xDocrequest["employer"].(string) != "" {
		whereemployer = fmt.Sprintf(`profile.employercontrol = '%s' AND`, xDocrequest["employer"])
	}

	whereemployertitle := ""
	if xDocrequest["employertitle"] != nil && xDocrequest["employertitle"].(string) != "" {
		whereemployertitle = fmt.Sprintf(`lower(employer.title) like lower('%%%s%%') AND`, xDocrequest["employertitle"])
	}

	wherecategory := ""
	if xDocrequest["category"] != nil && xDocrequest["category"].(string) != "" {
		wherecategory = fmt.Sprintf(`category.control = '%s' AND`, xDocrequest["category"])
	}

	wheresubcategory := ""
	if xDocrequest["subcategory"] != nil && xDocrequest["subcategory"].(string) != "" {
		wheresubcategory = fmt.Sprintf(`profile.subcategorycontrol = '%s' AND`, xDocrequest["subcategory"])
	}

	sqlFields := `
					profile.control as control, profile.workflow as workflow, profile.status as status,
					profile.employercontrol as employercontrol, profile.title as profiletitle,
					profile.firstname as firstname, profile.lastname as lastname, 
					profile.code as profilecode, profile.email as profileemail,
					scheme.control as schemecontrol, scheme.title as schemetitle, 
					subscription.startdate, subscription.expirydate
					
				`

	sqlOrderByLimit := `order by profile.firstname, profile.lastname limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(distinct profile.control) as paginationtotal	`
	}

	sql := `select %s

			from profile
            
            left join subscription on subscription.membercontrol = profile.control
			left join scheme on scheme.control = subscription.schemecontrol

            where %s %s %s %s %s %s %s %s %s %s  profile.company != 'Yes' and profile.control != '' ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, wherefirstname, wherelastname, whereemail, wherecode,
		whereworkflow, wherescheme, whereemployer, whereemployertitle, wherecategory, wheresubcategory, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Profile) Delete(xDocrequest map[string]interface{}, curdb Database) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `delete from %s where %s = '%s'`

	sql = fmt.Sprintf(sql, this.GetName(), searchfield, searchvalue)
	curdb.Query(sql)
}
