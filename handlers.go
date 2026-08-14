package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"

	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Notification struct {
	CreatedAt       time.Time
	Subj            string
	Message         string
	StudentFullName string
	TeacherFullName string
	ParentId        int
}

type AdminData struct {
	DistrictName string
	AddedSchools []SchoolStat
}
type TeacherLoginData struct {
	District     string
	LoginTeacher Teacher
	Students     []Student
}

type TeacherLoginError struct {
	Username string
	Error    string
	District string
}

type Student struct {
	Id             int
	First          string
	Last           string
	ParentFullName string
	ParentId       int
}

type Student2 struct {
	Id         int
	First      string
	Last       string
	SchoolName string
	SchoolId   int
}

type Teacher struct {
	Id       int
	First    string
	Last     string
	Username sql.NullString
	Grade    sql.NullString
	SchoolId int
}

type School struct {
	ID   string
	Name string
}

type SchoolStat struct {
	Name  string
	Id    int
	Count int
}

type MessageHistory struct {
	ID              int
	Subject         string
	StudentFullName string
	CreatedAt       time.Time
}

type MessageById struct {
	Message string
}
type RegisterPageData struct {
	First    string
	Last     string
	Error    string
	District string
}

type ParentDashData struct {
	ParentId       int
	ParentFullName string
	Notis          []Notification
}

type StudentU struct {
	Id     int
	Name   string
	Sid    int
	OldTid int
}

type AssignedStudentTeacher struct {
	Sid      int            //student id
	Sname    string         // student name
	Schoolid int            // student's schoolid
	Tid      int            // teacher's id
	Tname    string         //teacher's name
	Tgrade   sql.NullString //teacher's grade

}

type Student_Teacher struct {
	Name string
	Id   int
}

func saveTeacherPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("tp")), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error hashing")
		} else {
			err = os.WriteFile("datat", hashedPassword, 0600)
			if err != nil {
				fmt.Println("Error writing")
			}

			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

		}

	}
}

func saveAdminPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("ap")), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error hashing")
		} else {
			err = os.WriteFile("data", hashedPassword, 0600)
			if err != nil {
				fmt.Println("Error writing")
			}

			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

		}

	}
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	if r.Method == http.MethodPost { // User clicks Log In
		allowed := false
		hashpass, _ := os.ReadFile("data")
		entered_password := r.FormValue("password")

		if len(hashpass) == 0 {
			if entered_password == "orange" {
				allowed = true
			}
		} else { // comparing
			err := bcrypt.CompareHashAndPassword(hashpass, []byte(entered_password))
			if err == nil {
				//match
				allowed = true
			}
		}

		if allowed == false {
			//doesn't match
			http.Redirect(w, r, "/admin", http.StatusSeeOther)

		} else { //match
			//session, _ := store.Get(r, "session-name")
			session.Values["isAdmin"] = true

			err := session.Save(r, w)
			if err != nil {
				fmt.Println("Error saving session")
			}

			http.Redirect(w, r, "/admin/tasks", http.StatusSeeOther)

		}

	} else { //Presenting the form to user
		val, ok := session.Values["isAdmin"].(bool)
		if ok && val == true {
			http.Redirect(w, r, "/admin/tasks", http.StatusSeeOther)
			return
		}
		tmpl, err := template.ParseFiles("admin/login.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	}

}

func byschool_teachersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("public/byschool_teachers.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		db, err := GetCon()
		if err == nil {
			defer db.Close()
			schoolid := r.URL.Query().Get("schoolid")
			query := "SELECT id, firstname, lastname, schoolid FROM teachers Where schoolid = ?"

			rows, err := db.Query(query, schoolid)

			if err != nil {

				return
			}

			defer rows.Close()
			var teachers []Teacher

			for rows.Next() {
				var d Teacher
				rows.Scan(&d.Id, &d.First, &d.Last, &d.SchoolId)
				teachers = append(teachers, d)
			}

			tmpl.Execute(w, teachers)
		} else {
			w.Write([]byte("Can not get teachers"))
		}

	}
}

func addStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("public/hello.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, loadSchool())

	} else { // Post
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}
		defer db.Close()
		sfirstname := strings.TrimSpace(r.FormValue("firstname"))
		slastname := strings.TrimSpace(r.FormValue("lastname"))
		parentfullname := r.FormValue("parentfullname")
		parentid := r.FormValue("parentid")
		schoolid := r.FormValue("schoolid")

		_, err = db.Exec("INSERT INTO students (firstname, lastname, parent_fullname, parentid, schoolid) values (?, ?,?, ?,?)", sfirstname, slastname, parentfullname, parentid, schoolid)

		if err != nil {
			// MySQL duplicate error code = 1062
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Student already exists for this parent"))
				return
			}
			fmt.Println(err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		} //

		w.WriteHeader(http.StatusOK)

	}

}

func loadSchool() []School {
	db, err := GetCon()

	if err == nil { //no error

		defer db.Close()
		query := "SELECT id, name from schools"

		rows, err := db.Query(query)
		if err != nil {
			return nil
		}

		var schools []School
		for rows.Next() {
			var s School
			rows.Scan(&s.ID, &s.Name)
			schools = append(schools, s)
		}

		return schools
	}
	return nil
}

func AdminTasksHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")

	if err != nil {
		fmt.Println(err)
		return
	}

	val, ok := session.Values["isAdmin"].(bool)

	if !ok {

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return

	} else {

		if val == true {
			tmpl, err := template.ParseFiles("admin/tasks.html")
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
			districtbytes, err := os.ReadFile("setting.txt")
			districtstr := strings.TrimSpace(string(districtbytes))

			db, err := GetCon()

			if err == nil { //no error

				defer db.Close()
				schoolq := "select schools.name, schools.id, count(*) as teacher_count from schools join teachers on schools.id = teachers.schoolid group by schools.name, schools.id"
				rows, err := db.Query(schoolq)

				if err != nil {
					fmt.Println(err)

					admindata := AdminData{
						DistrictName: districtstr,
						AddedSchools: []SchoolStat{},
					}
					tmpl.Execute(w, admindata)

				} else {
					var schoolcounts []SchoolStat
					for rows.Next() {
						var d SchoolStat
						rows.Scan(&d.Name, &d.Id, &d.Count)
						schoolcounts = append(schoolcounts, d)
					}
					admindata := AdminData{
						DistrictName: districtstr,
						AddedSchools: schoolcounts,
					}

					tmpl.Execute(w, admindata)

				}

				defer rows.Close()
			} //no error

		}

	}

}

func CreateDistrictHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("dname")
	err := os.WriteFile("setting.txt", []byte(name), 0644)
	if err != nil {
		fmt.Println("Error storing district name")
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func createSchoolHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}

		schName := strings.TrimSpace(r.FormValue("schName"))
		// Insert into DB
		result, err := db.Exec("INSERT INTO schools (name) VALUES (?)", schName)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		lastid, err := result.LastInsertId()
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", lastid)
		}

	}

}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}
		tfirstname := strings.TrimSpace(r.FormValue("tfirst"))
		tlastname := strings.TrimSpace(r.FormValue("tlast"))
		tschoolid := r.FormValue("tschoolid")

		_, err = db.Exec("INSERT INTO teachers (firstname, lastname, schoolid) values (?, ?, ?)", tfirstname, tlastname, tschoolid)
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/admin/tasks", http.StatusSeeOther)

	}

}

func MessageByTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		db, err := GetCon()
		if err != nil {
			w.Write([]byte("Error connecting to database"))
			return
		}
		defer db.Close()

		pid := r.FormValue("parentid")
		tid := r.FormValue("teacherid")
		subj := r.FormValue("subj")
		msg := r.FormValue("message")
		sfame := r.FormValue("sfame")
		tfame := r.FormValue("teacher_fullname")

		_, err = db.Exec("INSERT INTO notifications(subject, message, student_fullname, teacher_fullname, parentid, teacherid) VALUES (?,?,?, ?, ?, ?)", subj, msg, sfame, tfame, pid, tid)
		if err != nil {
			fmt.Println("Error creating message", err)
			return
		}
		http.Redirect(w, r, "/teachers", http.StatusSeeOther)
	}

}

func updateGradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tid := r.URL.Query().Get("tid")
		tmpl, _ := template.ParseFiles("templates/teacher_grade.html")

		tmpl.Execute(w, tid)
	}

	if r.Method == http.MethodPost {
		ftid := r.FormValue("tid")
		grade := strings.TrimSpace(r.FormValue("grade"))

		db, err := GetCon()
		if err != nil {
			w.Write([]byte("Error connecting to database"))
			return
		}
		defer db.Close()

		query := "UPDATE teachers SET grade=? WHERE id=?"

		_, err = db.Exec(query, grade, ftid)
		if err != nil {
			w.Write([]byte("Error"))
		} else {
			w.Write([]byte("OK"))
		}
	}
}
func changableTeachersHandler(w http.ResponseWriter, r *http.Request) {
	schoolid := r.URL.Query().Get("schoolid")
	studentid := r.URL.Query().Get("stuid")

	db, err := GetCon()
	if err != nil {
		w.Write([]byte("Error connecting to database"))
		return
	}
	defer db.Close()

	q := "select teachers.id, teachers.firstname,  teachers.lastname, teachers.grade, teachers.schoolid from teachers where teachers.schoolid=? AND teachers.id not in (select st.tid from students_teachers st where st.sid=?)"
	rows, err := db.Query(q, schoolid, studentid)
	if err == nil {
		var tu []Teacher
		defer rows.Close()
		for rows.Next() {
			var d Teacher
			rows.Scan(&d.Id, &d.First, &d.Last, &d.Grade, &d.SchoolId)
			tu = append(tu, d)
		}
		tmpl, _ := template.ParseFiles("templates/changable_teachers.html")

		tmpl.Execute(w, tu)

	} else {
		fmt.Println(err)
	}

}

func changeTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := GetCon()
	if err != nil {
		w.Write([]byte("Error connecting to database"))
		return
	}
	defer db.Close()

	if r.Method == http.MethodGet {
		pid := r.URL.Query().Get("pid")
		q := "select s.id, concat(s.firstname,' ', s.lastname) , s.schoolid, st.tid,  concat( t.firstname, ' ',  t.lastname), t.grade from students s join students_teachers st on s.id = st.sid join teachers t on t.id= st.tid WHERE s.parentid=?"
		rows, err := db.Query(q, pid)

		if err == nil {
			var sus []AssignedStudentTeacher
			defer rows.Close()
			for rows.Next() {
				var d AssignedStudentTeacher
				rows.Scan(&d.Sid, &d.Sname, &d.Schoolid, &d.Tid, &d.Tname, &d.Tgrade)
				sus = append(sus, d)
			}
			tmpl, _ := template.ParseFiles("templates/change_teachers.html")

			tmpl.Execute(w, sus)

		}

	}

	if r.Method == http.MethodPost {
		studentid := r.FormValue("ssid")
		newteacherid := r.FormValue("ttid")
		oldteacherid := r.FormValue("otid")

		query := "UPDATE students_teachers SET tid=? WHERE sid=? AND tid=?"

		_, err = db.Exec(query, newteacherid, studentid, oldteacherid)
		if err != nil {

			fmt.Println(err)
			var mysqlErr *mysql.MySQLError

			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				w.Write([]byte("Already added the teacher."))
				return
			}
			w.Write([]byte("Error"))

		} else {
			w.Write([]byte("OK"))
		}
	}
}

func specifyTeacherHandler(w http.ResponseWriter, r *http.Request) {

	db, err := GetCon()
	if err != nil {
		w.Write([]byte("Error connecting to database"))
		return
	}
	defer db.Close()
	if r.Method == http.MethodGet {
		pid := r.URL.Query().Get("pid")

		q := "select s.id, concat(s.firstname,' ', s.lastname) , s.schoolid from students s join users u on s.parentid = u.id where u.id = ?;"
		rows, err := db.Query(q, pid)

		if err == nil {
			var sus []StudentU
			defer rows.Close()
			for rows.Next() {
				var d StudentU
				rows.Scan(&d.Id, &d.Name, &d.Sid)
				sus = append(sus, d)
			}
			tmpl, _ := template.ParseFiles("templates/specify_teacher.html")
			tmpl.Execute(w, sus)

		}
	}

	if r.Method == http.MethodPost {
		studentid := r.FormValue("stid")
		teacherid := r.FormValue("tid")
		schoolid := r.URL.Query().Get("schoolid")
		query := "INSERT INTO students_teachers (sid, tid, schoolid) VALUES (?, ?, ?)"

		_, err = db.Exec(query, studentid, teacherid, schoolid)
		if err != nil {

			var mysqlErr *mysql.MySQLError

			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				w.Write([]byte("Already added the teacher."))
				return
			}
			w.Write([]byte("Error"))

		} else {
			w.Write([]byte("OK"))
		}
	}

}

func teacherUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db, err := GetCon()
		if err != nil {
			w.Write([]byte("Error connecting to database"))
			return
		}
		defer db.Close()
		uname := r.FormValue("tname")
		grade := r.FormValue("tgrade")
		id := r.FormValue("tid")
		query := "UPDATE teachers SET  username= ?, grade=?  WHERE id = ?"

		_, err = db.Exec(query, uname, grade, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))

	}
}

func teacherSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		db, err := GetCon()
		if err != nil {
			w.Write([]byte("Error connecting to database"))
			return
		}
		defer db.Close()
		tsid := r.FormValue(("tsid"))
		tfn := r.FormValue("tfn")
		tln := r.FormValue("tln")

		tquery := "SELECT id, firstname, lastname, username, grade, schoolid FROM teachers WHERE schoolid = ? AND firstname = ? AND lastname = ?"
		rows, err := db.Query(tquery, tsid, tfn, tln)

		if err == nil {
			var teachers []Teacher
			//Load Teacher Dashboard
			for rows.Next() {
				var d Teacher
				rows.Scan(&d.Id, &d.First, &d.Last, &d.Username, &d.Grade, &d.SchoolId)
				teachers = append(teachers, d)
			}
			if len(teachers) == 0 {
				w.Write([]byte("Not Found"))
				return
			}

			tmpl := template.Must(template.ParseFiles("templates/teachers_search.html"))
			err := tmpl.Execute(w, teachers)
			if err != nil {
				log.Printf("template execute error: %v", err)
			}
		} else {
			fmt.Println(err)
		}

	}
}

func usernameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/username.html"))
	tmpl.Execute(w, loadSchool())
}

func TeachersLoginHandler(w http.ResponseWriter, r *http.Request) {
	allowed := false
	district, _ := os.ReadFile("setting.txt")
	dis := string(district)
	tmpl := template.Must(template.ParseFiles("templates/teachers_login.html"))

	var le TeacherLoginError
	session, _ := store.Get(r, "session-name")

	if r.Method == http.MethodPost {
		hashpass, _ := os.ReadFile("datat")
		entered_password := r.FormValue("password")
		un := r.FormValue("username")
		if len(hashpass) == 0 { // datat does not exist, or exist but empty
			if entered_password == "chalk" {
				//match default password for teacher
				allowed = true
			} else {
				//Default password doesn't match
				le = TeacherLoginError{Username: un, Error: "Password does not match", District: dis}
				tmpl.Execute(w, le)

			}
		} else { // comparing stored passowrd and entered password

			err := bcrypt.CompareHashAndPassword(hashpass, []byte(entered_password))
			if err == nil {
				//password match
				allowed = true

			} else {
				//password doesn't match
				le = TeacherLoginError{Username: un, Error: "Password does not match", District: dis}
				tmpl.Execute(w, le)

			}
		}
		if allowed == true {
			session.Values["un"] = un

			session.Save(r, w)
			loadTeacherDashboard(w, tmpl, un, dis)
		}

	} else { //get method

		uname, ok := session.Values["un"].(string)

		if ok {
			loadTeacherDashboard(w, tmpl, uname, dis)
		} else {
			le = TeacherLoginError{Username: "", Error: "", District: dis}
			tmpl.Execute(w, le)
		}

	}
}

func loadTeacherDashboard(w http.ResponseWriter, tmpl *template.Template, username string, dist string) {

	db, err := GetCon()
	if err == nil {
		defer db.Close()

		row := db.QueryRow("SELECT id, firstname, lastname, username, grade, schoolid FROM teachers WHERE username =?", username)

		var teacherid, schoolid int
		var firstName, lastName, loginame, grade string

		err := row.Scan(&teacherid, &firstName, &lastName, &loginame, &grade, &schoolid)
		if err == nil { //found teacher
			login_teacher := Teacher{
				Id:       teacherid,
				First:    firstName,
				Last:     lastName,
				SchoolId: schoolid,
			}
			//Load students for the teacher
			student_query := "SELECT s.id, s.firstname, s.lastname, s.parent_fullname, s.parentid FROM students s join students_teachers ss on s.id = ss.sid WHERE ss.tid=?"
			rows, err := db.Query(student_query, teacherid)

			if err == nil {
				var Query_Students []Student
				//Load Teacher Dashboard
				for rows.Next() {
					var d Student
					rows.Scan(&d.Id, &d.First, &d.Last, &d.ParentFullName, &d.ParentId)
					Query_Students = append(Query_Students, d)
				}

				Tld := TeacherLoginData{
					District:     string(dist),
					LoginTeacher: login_teacher,
					Students:     Query_Students,
				}
				tdb_tmpl := template.Must(template.ParseFiles("templates/teachers_dashboard.html"))
				tdb_tmpl.Execute(w, Tld)

			} else {
				fmt.Println("Error querying students", err)
			}

			defer rows.Close()

		} else {

			le := TeacherLoginError{Username: "", Error: username + " is not found", District: dist}
			tmpl.Execute(w, le)
			return
		}

	} else {
		fmt.Println("Error connecting to database in TeachersLoginHandler()")
	}

}
func getMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tid := r.URL.Query().Get("tid")
	db, err := GetCon()
	if err == nil {
		defer db.Close()
		msg_hisq := "SELECT id, subject, student_fullname, created_at FROM notifications  WHERE teacherid = ? order by created_at DESC"

		msgrows, msgerr := db.Query(msg_hisq, tid)
		var MsgHis []MessageHistory
		if msgerr == nil {
			for msgrows.Next() {
				var h MessageHistory
				msgrows.Scan(&h.ID, &h.Subject, &h.StudentFullName, &h.CreatedAt)
				MsgHis = append(MsgHis, h)

			}

			tmpl := template.Must(template.ParseFiles("public/msg_his.html"))
			tmpl.Execute(w, MsgHis)
		}
		defer msgrows.Close()
	}

}

func messageDeleteHandler(w http.ResponseWriter, r *http.Request) {
	did := r.URL.Query().Get("id")
	db, err := GetCon()
	if err == nil {

		qd := "DELETE FROM notifications WHERE id = ?"
		db.Exec(qd, did)

	}
	defer db.Close()

}

func messageEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		did := r.URL.Query().Get("id")
		msg := r.FormValue("a_" + did)

		db, err := GetCon()
		if err == nil {
			qd := "UPDATE notifications SET message = ? WHERE id = ?"
			db.Exec(qd, msg, did)

		}
		defer db.Close()
	}
}

func messageByIdHandler(w http.ResponseWriter, r *http.Request) {
	msgid := r.URL.Query().Get("mid")

	db, err := GetCon()
	if err == nil {
		defer db.Close()
		q := "SELECT message from notifications WHERE id = ?"
		rows, err := db.Query(q, msgid)
		if err == nil {
			var Msgs MessageById
			for rows.Next() {
				rows.Scan(&Msgs.Message)

			}

			tmpl := template.Must(template.ParseFiles("public/msgbyid.html"))
			tmpl.Execute(w, Msgs)

		}
		defer rows.Close()
	}

}
func uploadSchoolsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}

		defer db.Close()
		file, _, err := r.FormFile("schUploadFile")

		if err != nil {
			http.Error(w, "Invalid file upload", http.StatusBadRequest)
			return
		}

		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll() // Read all rows into a slice of slices
		if err != nil {
			http.Error(w, "Error parsing CSV", http.StatusInternalServerError)
			return
		}
		query := "INSERT INTO schools (name) VALUES "
		values := []interface{}{}

		for i, row := range records {
			query += "(?)"
			if i < len(records)-1 {
				query += ","
			}

			values = append(values, strings.TrimSpace(string(row[0])))

		}
		_, err = db.Exec(query, values...)

		if err == nil {

			http.Redirect(w, r, "/admin/tasks", http.StatusSeeOther)
		}

	}
}

func uploadTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}

		defer db.Close()
		schoolid := r.FormValue("tschoolid2")
		file, _, err := r.FormFile("csvfile")

		if err != nil {
			http.Error(w, "Invalid file upload", http.StatusBadRequest)
			return
		}

		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll() // Read all rows into a slice of slices
		if err != nil {
			http.Error(w, "Error parsing CSV", http.StatusInternalServerError)
			return
		}
		query := "INSERT INTO teachers (firstname, lastname, schoolid) VALUES "
		values := []interface{}{}

		for i, row := range records {
			query += "(?, ?, ?)"
			if i < len(records)-1 {
				query += ","
			}
			ischoolid, _ := strconv.Atoi(schoolid)

			values = append(values, strings.TrimSpace(string(row[0])), strings.TrimSpace(string(row[1])), ischoolid)

		}
		_, err = db.Exec(query, values...)

		if err == nil {

			http.Redirect(w, r, "/admin/tasks", http.StatusSeeOther)
		}

	}

}
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	defer conn.Close()
	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Render the home page with a simple greeting
	district, _ := os.ReadFile("setting.txt")
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, string(district))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/register.html"))
	district, _ := os.ReadFile("setting.txt")

	switch r.Method {
	case http.MethodGet:
		tmpl.Execute(w, RegisterPageData{
			First:    "",
			Last:     "",
			Error:    "",
			District: string(district),
		})
		return
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		firstname := strings.TrimSpace(r.FormValue("firstname"))
		lastname := strings.TrimSpace(r.FormValue("lastname"))

		// Hash password hash,
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		db, err := GetCon()
		if err != nil {
			fmt.Println("Error connecting to database")
			return
		}
		// Insert user
		_, err = db.Exec(` INSERT INTO users (username, firstname, lastname, password_hash) VALUES (?, ?, ?, ?) `, username, firstname, lastname, string(hash))
		if err != nil { // MySQL duplicate entry error
			if isDuplicateError(err) {
				tmpl.Execute(w,
					RegisterPageData{
						First:    firstname,
						Last:     lastname,
						Error:    username + " already taken",
						District: string(district)})
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		// Success → redirect
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}

// Login handler
func parentHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	val, _ := session.Values["isUserIn"].(bool)

	if val == true {
		//redisplay the parent_dash.html
		pfname := r.URL.Query().Get("pfn")
		pids := r.URL.Query().Get("pid")
		pid, err := strconv.Atoi(pids)
		if err != nil {
			return
		}

		loadParentPage(w, pid, pfname, "public/parent_fetch.html")
	}

}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		db, err := GetCon()

		if err != nil {
			http.Error(w, "Error connecting", http.StatusUnauthorized)
		}

		defer db.Close()

		query := "SELECT id, username, CONCAT(firstname, ' ', lastname) as parentfullname, password_hash FROM users WHERE username=?"
		row := db.QueryRow(query, username)

		var id int
		var uname, hpass, parentfullname string

		err = row.Scan(&id, &uname, &parentfullname, &hpass)

		if err == nil {

			//compare password
			err_comp := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(password))

			if err_comp == nil {

				session.Values["isUserIn"] = true
				session.Values["pid"] = id
				session.Values["pfn"] = parentfullname

				errs := session.Save(r, w)
				if errs != nil {
					fmt.Println("Error saving session")
				}
				loadParentPage(w, id, parentfullname, "public/parent_dash.html")
				return
			}

		}

	}

	// Show the login form
	district, _ := os.ReadFile("setting.txt")
	tmpl := template.Must(template.ParseFiles("templates/login.html"))

	val, _ := session.Values["isUserIn"].(bool)

	if val != true {
		tmpl.Execute(w, string(district))
	} else {
		if padi, ok := session.Values["pid"].(int); ok {
			pfname := session.Values["pfn"].(string)
			loadParentPage(w, padi, pfname, "public/parent_dash.html")

		} else {

			w.Write([]byte("Error"))
		}

	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/", http.StatusFound)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		log.Println("recv:", string(msg))

		err = conn.WriteMessage(msgType, msg) // echo
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func getSchoolsHandler(w http.ResponseWriter, r *http.Request) {

	db, err := GetCon()

	if err != nil {

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT id, name FROM schools"

	rows, err := db.Query(query)

	if err != nil {
		// Handle specific errors (e.g., duplicate ID)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var schools []School

	for rows.Next() {
		var d School
		rows.Scan(&d.ID, &d.Name)
		schools = append(schools, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schools)
}

func loadParentPage(ww http.ResponseWriter, pd int, pn string, tmpl string) {
	db, err := GetCon()
	if err != nil {
		http.Error(ww, "Error connecting", http.StatusUnauthorized)
	}

	defer db.Close()

	pq := "SELECT created_at, subject, message, student_fullname, teacher_fullname,  parentid FROM notifications WHERE parentid = ?"
	rows, err := db.Query(pq, pd)
	var Notifications []Notification
	if err == nil {

		for rows.Next() {
			var an Notification
			rows.Scan(&an.CreatedAt, &an.Subj, &an.Message, &an.StudentFullName, &an.TeacherFullName, &an.ParentId)
			Notifications = append(Notifications, an)
		}

	} else {
		fmt.Println("Error in select notifications", err)
	}

	defer rows.Close()
	ptmpl := template.Must(template.ParseFiles(tmpl))
	pdata := ParentDashData{
		ParentId:       pd,
		ParentFullName: pn,
		Notis:          Notifications,
	}

	ptmpl.Execute(ww, pdata)
}

func deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	_, ok := session.Values["un"].(string)
	ptmpl := template.Must(template.ParseFiles("templates/delete_student.html"))

	if !ok {
		w.Write([]byte("Teacher should sign in"))
		return
	}

	db, err := GetCon()
	if err != nil {
		http.Error(w, "Error connecting", http.StatusUnauthorized)
	}

	defer db.Close()
	if r.Method == http.MethodGet {
		tid := r.URL.Query().Get("tid")
		pq := "select concat(s.firstname,' ',  s.lastname) as name , st.sid from students s join students_teachers st on s.id=st.sid where st.tid = ?"
		rows, err := db.Query(pq, tid)
		var Students_Teachers []Student_Teacher
		if err == nil {
			for rows.Next() {
				var ast Student_Teacher
				rows.Scan(&ast.Name, &ast.Id)
				Students_Teachers = append(Students_Teachers, ast)
			}

		} else {
			fmt.Println("Error in selecting data: ", err)

		}

		defer rows.Close()

		ptmpl.Execute(w, Students_Teachers)

	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		ids := r.Form["sids"]

		if len(ids) == 0 {
			http.Redirect(w, r, "/teachers", http.StatusSeeOther)
			return
		}
		q := "DELETE FROM students_teachers WHERE sid in " + "(" + strings.Join(ids, ",") + ")"
		_, err = db.Exec(q)

		if err != nil {
			http.Error(w, "Delete failed. Please do go /teachers", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/teachers", http.StatusSeeOther)

	}
}

func changeSchoolHandler(w http.ResponseWriter, r *http.Request) {
	parentid := r.URL.Query().Get("pid")

	db, err := GetCon()
	if err != nil {
		http.Error(w, "Error connecting", http.StatusUnauthorized)
	}

	defer db.Close()
	if r.Method == http.MethodGet {
		q1 := "select s.id, s.firstname, s.lastname, l.name, l.id from students s join schools l on s.schoolid=l.id where s.parentid=?"
		rows, err := db.Query(q1, parentid)
		var st2 []Student2
		if err == nil {

			for rows.Next() {
				var s2 Student2
				rows.Scan(&s2.Id, &s2.First, &s2.Last, &s2.SchoolName, &s2.SchoolId)
				st2 = append(st2, s2)
			}

		} else {
			fmt.Println(err)

		}
		defer rows.Close()

		p2 := template.Must(template.ParseFiles("templates/change_school.html"))

		p2.Execute(w, st2)

	}
	if r.Method == http.MethodPost {
		pp := r.URL.Query().Get("pid")
		studentid := r.FormValue("stutid")
		schoolid := r.FormValue("scuid")
		oldschoolid := r.URL.Query().Get("oldschoolid")

		ex := "update students set schoolid=? where id=? and parentid=?"
		_, error := db.Exec(ex, schoolid, studentid, pp)
		if error == nil {
			ex2 := "Delete from students_teachers where sid = ? and schoolid=?"
			_, error2 := db.Exec(ex2, studentid, oldschoolid)

			if error2 == nil {
				w.Write([]byte("OK"))
			}
		}

	}
}

func changeable_SchoolHandler(w http.ResponseWriter, r *http.Request) {
	studentid := r.URL.Query().Get("sid")
	db, err := GetCon()
	if err != nil {
		http.Error(w, "Error connecting", http.StatusUnauthorized)
	}

	defer db.Close()
	q1 := "select id, name from schools where id not in (Select schoolid from students where id=?)"
	rows, error := db.Query(q1, studentid)
	var schools []School
	if error == nil {

		for rows.Next() {
			var s2 School
			rows.Scan(&s2.ID, &s2.Name)
			schools = append(schools, s2)
		}
		p2 := template.Must(template.ParseFiles("templates/changeable_school.html"))

		p2.Execute(w, schools)
	}

	defer rows.Close()
}
