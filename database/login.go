package database

import (
	"encoding/base64"
	"fmt"
)

type Login struct {
	Crud
}

func (this *Login) GetName() string {
	return "login"
}

func (this *Login) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Login) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["image"] = ""
	fields["role"] = ""
	fields["username"] = ""
	fields["password"] = ""
	fields["mobilecode"] = ""
	fields["mobile"] = ""
	fields["email"] = ""
	fields["profilecontrol"] = ""

	return fields
}

func (this *Login) EncryptPassword(cString string, curdb Database) string {

	byteEncrypted := curdb.Encrypt([]byte(cString))
	return base64.StdEncoding.EncodeToString(byteEncrypted)
	// return base64.StdEncoding.EncodeToString([]byte(cString))
}

func (this *Login) VerifyLogin(lFrontend bool, cUsername string, cPassword string, curdb Database) (mapRes map[string]interface{}) {

	sFrontend := `AND login.role != 'member'`
	if lFrontend {
		sFrontend = ``
	}

	sql := `select distinct	login.code as code, login.control as control, login.title as title, 
					login.description as description, login.workflow as workflow,
					login.createdate as createdate, login.createdby as createdby, 
					login.updatedate as updatedate, login.updatedby as updatedby,
					
					login.role as role, login.username as username, login.password as password, 
					login.email as email, login.mobile as mobile, login.mobilecode as mobilecode, login.image as image, 

					login.profilecontrol as profilecontrol, profile.title as profiletitle, profile.merchantpin as merchantpin,
					profile.firstname as profilefirstname, profile.lastname as profilelastname,
					profile.image as profileimage, profile.employercontrol as employercontrol, employer.title as employertitle,
					profile.workflow as profileworkflow

			from login left join profile on login.profilecontrol = profile.control
			left join profile as employer on profile.employercontrol = employer.control

			where login.username = '%s' AND login.password = '%s' AND login.workflow = 'active'  %s
			`
	sql = fmt.Sprintf(sql, cUsername, cPassword, sFrontend)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Login) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["code"] = "root"
	xDoc["role"] = "admin"
	xDoc["title"] = "Root"
	xDoc["username"] = "root@localhost.com"
	xDoc["password"] = this.EncryptPassword("toor", curdb)
	xDoc["workflow"] = "active"
	xDoc["profilecontrol"] = "1.00000000001"
	this.Create("root", xDoc, curdb)
}

func (this *Login) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Login) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Login) ReadUser(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select	login.code as code, login.control as control, login.title as title, 
					login.description as description, login.workflow as workflow,
					login.createdate as createdate, login.createdby as createdby, 
					login.updatedate as updatedate, login.updatedby as updatedby,
					
					login.role as role, login.username as username, login.password as password, 
					login.email as email, login.mobile as mobile, login.mobilecode as mobilecode, login.image as image, 

					login.profilecontrol as profilecontrol, profile.title as profiletitle, profile.merchantpin as merchantpin,
					profile.firstname as profilefirstname, profile.lastname as profilelastname, profile.email as profileemail
					profile.image as profileimage, profile.employercontrol as employercontrol

			from login left join profile on login.profilecontrol = profile.control

				where login.%s = '%s' order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Login) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select  distinct login.username as username,
					login.workflow as accountstatus,  login.role as usertype, login.role as role,  login.control as logincontrol,
					login.code as pin, login.password as password, login.profilecontrol as profilecontrol,

					profile.code as code, profile.control as control, profile.title as title, profile.merchantpin as merchantpin,
					profile.description as description, profile.workflow as workflow, 
					profile.createdate as createdate, profile.createdby as createdby, 
					profile.updatedate as updatedate, profile.updatedby as updatedby,

					profile.firstname as firstname, profile.lastname as lastname, profile.referrer as referrer, 
					profile.gender as gender, profile.dob as dob, profile.nationality as nationality, 

					profile.website as website, profile.phonecode as phonecode, profile.phone as phone, profile.email as email, 
					profile.emailsecondary as emailsecondary, profile.emailalternate as emailalternate,

					profile.member as member, profile.merchant as merchant, profile.employer as employer,
					profile.image as image, profile.keywords as keywords, profile.merchantterms as merchantterms, 
					profile.merchantbenefits as merchantbenefits, 

					profile.employercontrol as employercontrol, employer.title as employertitle,
					industry.title as industrytitle, profile.industrycontrol as industrycontrol,
					profile.subindustrycontrol as subindustrycontrol,
					
					profile.categorycontrol as categorycontrol, category.title as categorytitle,
                    profile.subcategorycontrol as subcategorycontrol, subcategory.title as subcategorytitle


			from profile
            
            left join industry on profile.industrycontrol = industry.control 
                                               
            left join category on profile.categorycontrol = category.control
            left join category as subcategory on profile.subcategorycontrol = subcategory.control
                       
            left join profile as employer on profile.employercontrol = employer.control 
            left join login on login.profilecontrol = profile.control
            
            
            where profile.%s = '%s' 
					
			order by firstname, lastname, title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Login) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

	whererole := ""
	if xDocrequest["role"] != nil && xDocrequest["role"].(string) != "" {
		whererole = fmt.Sprintf(`lower(login.role) like lower('%%%s%%') AND`, xDocrequest["role"])
	}

	whereaccountstatus := ""
	if xDocrequest["accountstatus"] != nil && xDocrequest["accountstatus"].(string) != "" {
		whereaccountstatus = fmt.Sprintf(`login.workflow = '%s' AND`, xDocrequest["accountstatus"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`profile.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`lower(profile.title) like lower('%%%s%%') AND`, xDocrequest["title"])
	}

	wheregender := ""
	if xDocrequest["gender"] != nil && xDocrequest["gender"].(string) != "" {
		wheregender = fmt.Sprintf(`lower(profile.gender) like lower('%%%s%%') AND`, xDocrequest["gender"])
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
		whereusername = fmt.Sprintf(`lower(login.username) like lower('%%%s%%') AND`, xDocrequest["username"])
	}

	whereexpirydate := ""
	if xDocrequest["expirydate"] != nil && xDocrequest["expirydate"].(string) != "" {
		whereexpirydate = fmt.Sprintf(`profile.expirydate like '%%%s%%' AND`, xDocrequest["expirydate"])
	}

	// whereprofilecontrol := ""
	// if xDocrequest["profilecontrol"] != nil && len(xDocrequest["profilecontrol"].([]string)) > 0 {
	// 	whereprofilecontrol = strings.Join(xDocrequest["profilecontrol"].([]string), `","`)
	// 	whereprofilecontrol = fmt.Sprintf(`profile.control in ("%s") AND`, whereprofilecontrol)
	// }

	whereprofilecontrol := ""
	if xDocrequest["profilecontrol"] != nil && xDocrequest["profilecontrol"].(string) != "" {
		whereprofilecontrol = fmt.Sprintf(`login.profilecontrol = '%s' AND`, xDocrequest["profilecontrol"])
	}

	whereprofile := ""
	if xDocrequest["member"] != nil && xDocrequest["member"].(string) != "" {
		whereprofile = fmt.Sprintf(`profile.member = '%s' AND`, xDocrequest["member"])
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`profile.merchant = '%s' AND`, xDocrequest["merchant"])
	}

	whereemployer := ""
	if xDocrequest["employer"] != nil && xDocrequest["employer"].(string) != "" {
		whereemployer = fmt.Sprintf(`profile.employer = '%s' AND`, xDocrequest["employer"])
	}

	whereemployertitle := ""
	if xDocrequest["employertitle"] != nil && xDocrequest["employertitle"].(string) != "" {
		whereemployertitle = fmt.Sprintf(`lower(employer.title) like lower('%%%s%%') AND`, xDocrequest["employertitle"])
	}

	whereemployercontrol := ""
	if xDocrequest["employercontrol"] != nil && xDocrequest["employercontrol"].(string) != "" {
		whereemployercontrol = fmt.Sprintf(`profile.employercontrol = '%s' AND`, xDocrequest["employercontrol"])
	}

	whereindustryname := ""
	if xDocrequest["industryname"] != nil && xDocrequest["industryname"].(string) != "" {
		whereindustryname = fmt.Sprintf(`lower(industry.title) like lower('%%%s%%') AND`, xDocrequest["industryname"])
	}

	whereindustry := ""
	if xDocrequest["industry"] != nil && xDocrequest["industry"].(string) != "" {
		whereindustry = fmt.Sprintf(`industry.control = '%s' AND`, xDocrequest["industry"])
	}

	wheresubindustry := ""
	if xDocrequest["subindustry"] != nil && xDocrequest["subindustry"].(string) != "" {
		wheresubindustry = fmt.Sprintf(`profile.subindustrycontrol = '%s' AND`, xDocrequest["subindustry"])
	}

	wherecategoryname := ""
	if xDocrequest["categoryname"] != nil && xDocrequest["categoryname"].(string) != "" {
		wherecategoryname = fmt.Sprintf(`lower(category.title) like lower('%%%s%%') AND`, xDocrequest["categoryname"])
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

	sqlFields := `
						distinct login.username as username,
						login.workflow as accountstatus,  login.role as usertype, login.role as role,  login.control as logincontrol,
						login.code as pin, login.password as password,  login.profilecontrol as profilecontrol,

						profile.code as code, profile.control as control, profile.title as title, profile.merchantpin as merchantpin, 
						profile.description as description, profile.workflow as workflow, 
						profile.createdate as createdate, profile.createdby as createdby, 
						profile.updatedate as updatedate, profile.updatedby as updatedby,

						profile.firstname as firstname, profile.lastname as lastname, profile.referrer as referrer, 
						profile.gender as gender, profile.dob as dob, profile.nationality as nationality, 

						profile.website as website, profile.phonecode as phonecode, profile.phone as phone, profile.email as email, 
						profile.emailsecondary as emailsecondary, profile.emailalternate as emailalternate,

						profile.member as member, profile.merchant as merchant, profile.employer as employer,
						profile.image as image, profile.keywords as keywords, profile.merchantterms as merchantterms, 
						profile.merchantbenefits as merchantbenefits, 

						profile.employercontrol as employercontrol, employer.title as employertitle,
						industry.title as industrytitle, profile.industrycontrol as industrycontrol,
						profile.subindustrycontrol as subindustrycontrol,
						
						profile.categorycontrol as categorycontrol, category.title as categorytitle,
						profile.subcategorycontrol as subcategorycontrol, subcategory.title as subcategorytitle
					
					`

	sqlOrderByLimit := `order by profile.firstname, profile.lastname, profile.title limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(distinct login.username) as paginationtotal	`
	}

	sql := `select %s

				from login

				left join profile on profile.control = login.profilecontrol
				left join industry on profile.industrycontrol = industry.control 
				
				left join category on profile.categorycontrol = category.control
				left join category as subcategory on profile.subcategorycontrol = subcategory.control
				
				left join profile as employer on profile.employercontrol = employer.control 
				
	            
	            where  %s %s %s %s %s %s %s %s %s %s %s %s %s %s  %s %s  %s  %s  %s %s  %s %s  %s %s 
				
				login.control != '' ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, wheregender, wherefirstname, wherelastname,
		wherereferrer, whereemail, whereexpirydate, whereemployer, whereemployertitle, whereprofile, wheremerchant, whereusername, whererole, whereaccountstatus,
		whereemployercontrol, whereindustry, whereindustryname, wheresubindustry, wherecategory, wherecategoryname,
		wheresubcategory, wherekeywords, whereprofilecontrol, limit, offset)

	mapRes, _ = curdb.Query(sql)
	return mapRes

}
