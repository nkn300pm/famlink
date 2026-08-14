
document.getElementById("csfrm").addEventListener('submit', async (e) => {
    e.preventDefault();
    f=e.target;
    fdt = new FormData(f);
    const purl=document.getElementById("csd");
    const sdent= document.getElementById("stutid");
    const sitem= sdent.options[sdent.selectedIndex];      


    try{
        const response = await fetch(purl.getAttribute("href") + "&oldschoolid="+ sitem.dataset.oldschoolid, {
            method: 'POST',
            body: fdt
        });
        const data=await response.text();
        if (data=="OK"){
            const newscu=document.getElementById("scuid");
            const sschool=newscu.options[newscu.selectedIndex];
            const newsname=sschool.text;
            alert("Successfully changed to " + newsname); 
            sitem.setAttribute("data-oldschoolid", sschool.value);               
        
                       
            const sname=sitem.text.split(" -- ")[0];
            sitem.text = sname + " -- " + newsname;           
            newscu.innerHTML="";  
            sdent.selectedIndex=0;
           
        }else{
            alert(data);
        }
    }catch {};
});

document.getElementById("stutid").addEventListener('change', (e) => {
    e.preventDefault();
    const slist=e.target;
   const so= slist.options[slist.selectedIndex];
   chid=so.value;

   try{
    fetch('/changable_school?sid=' + chid).then(resp => resp.text()).then(html => {
       
       document.getElementById('scuid').innerHTML = html;
    })


}catch(error) { alert(error)};


});
