var sstudent;
var schoolid;
document.getElementById("stid").addEventListener('change', (e) => {
    e.preventDefault();
    const list=e.target;
   const s= list.options[list.selectedIndex];
   schoolid=s.dataset.sid;
   sstudent = s.text;  
 
       try{
        fetch('/changableTeachers?schoolid=' + schoolid + '&stuid=' + s.value).then(resp => resp.text()).then(html => {
           shtml = '<option value="" disabled selected>Select a teacher</option>' + html;
           document.getElementById('tid').innerHTML = shtml;
        })


}catch(error) { alert(error)};

});

document.getElementById("sfrm").addEventListener('submit', async (e) => {
    e.preventDefault();
    fr=e.target;
    fd = new FormData(fr);
    try{
        const response = await fetch('/us?schoolid=' + schoolid, {
            method: 'POST',
            body: fd
        });
        const data=await response.text();
        if (data=="OK"){
            const t=fr.elements["tid"];
            document.getElementById("stat").innerHTML="Now " + sstudent + "'s teacher is " + t.options[t.selectedIndex].text + ". Next school year, if having new teacher; use Change Teacher menu";
        }else{
            alert(data);
        }
    }catch {};
});