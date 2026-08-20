var sname;
document.getElementById("ssid").addEventListener('change', (e) => {
    e.preventDefault();
    const list=e.target;
   const so= list.options[list.selectedIndex];
   scuid=so.dataset.schoolid;
   oldtid=so.dataset.oldtid;
  sname = so.text.split('--')[0];
  document.getElementById('otid').value=oldtid;

   try{
    fetch('/changableTeachers?schoolid=' + scuid + '&stuid=' + so.value).then(resp => resp.text()).then(html => {
       sl = '<option value="" disabled selected>Select a teacher</option>' + html;
       document.getElementById('ttid').innerHTML = sl;
    });

}catch(error) { alert(error)};

});

document.getElementById("cfrm").addEventListener('submit', async (e) => {
    e.preventDefault();
    f=e.target;
    fdt = new FormData(f);
    try{
        const response = await fetch('/ct', {
            method: 'POST',
            body: fdt
        });
        const data=await response.text();
        if (data=="OK"){
            const t=f.elements["ttid"];
            document.getElementById("box1").innerHTML=`<div class="alert alert-success d-flex align-items-center">
    <strong>Success!</strong>&nbsp;${"Now " + sname + "'s teacher is " + t.options[t.selectedIndex].text}</div>`;            
           
        }else{
            alert(data);
        }
    }catch {};
});