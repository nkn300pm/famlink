package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

var adminlogin = false

var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   28800,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/admin/teacher/add", addTeacherHandler)
	mux.HandleFunc("/admin/teachers/upload", uploadTeacherHandler)
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/echo", wsHandler)
	mux.HandleFunc("/admin", AdminLoginHandler)
	mux.HandleFunc("/admin/tasks", AdminTasksHandler)
	mux.HandleFunc("/parents", parentHandler)
	mux.HandleFunc("/teachers", TeachersLoginHandler)
	mux.HandleFunc("/messageById", messageByIdHandler)
	mux.HandleFunc("/message", MessageByTeacherHandler)
	mux.HandleFunc("/mdel", messageDeleteHandler)
	mux.HandleFunc("/medit", messageEditHandler)
	mux.HandleFunc("/mhis", getMessageHistoryHandler)
	mux.HandleFunc("/admin/district/add", CreateDistrictHandler)
	mux.HandleFunc("/admin/school/add", createSchoolHandler)
	mux.HandleFunc("/admin/school/upload", uploadSchoolsHandler)
	mux.HandleFunc("/admin/getschools", getSchoolsHandler)
	mux.HandleFunc("/admin/district/ap", saveAdminPasswordHandler)
	mux.HandleFunc("/admin/district/tp", saveTeacherPasswordHandler)
	mux.HandleFunc("/add/students", addStudentsHandler)
	mux.HandleFunc("/byschool_teachers", byschool_teachersHandler)
	mux.HandleFunc("/username", usernameHandler)
	mux.HandleFunc("/ts", teacherSearchHandler)
	mux.HandleFunc("/ut", teacherUpdateHandler)
	mux.HandleFunc("/us", specifyTeacherHandler)
	mux.HandleFunc("/ct", changeTeacherHandler)
	mux.HandleFunc("/changableTeachers", changableTeachersHandler)
	mux.HandleFunc("/udg", updateGradeHandler)
	mux.HandleFunc("/ds", deleteStudentHandler)
	mux.HandleFunc("/cs", changeSchoolHandler)
	mux.HandleFunc("/changable_school", changeable_SchoolHandler)

	server := &http.Server{
		Addr:              ":8085",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on", server.Addr)
	log.Fatal(server.Serve(ln))
}
