
document.querySelectorAll('.list-group-item').forEach(item => {
    item.addEventListener('click', function (e) {
        e.preventDefault();        
       
        if (item.getAttribute('v') === null) {
            iid=item.getAttribute('id').substring(1);                
        fetch('/messageById?mid=' + iid).then(res => res.text()).then(html => {            
        document.getElementById('e_' + iid).hidden=false;        
        const ta = document.getElementById('a_'+iid);
        ta.hidden=false;
        ta.value= html;  
         item.setAttribute('v', 'v');     
         document.getElementById('stat').innerHTML='';
       
        });

        }         
    });
});   

document.querySelectorAll('a[id^="d_"]').forEach(item => {
    item.addEventListener('click', async function(e) {
        e.preventDefault();
        did = item.getAttribute('id').substring(2);
        const furl="/mdel?id=" + did;

        try {
            const response = await fetch(furl);
            if (response.ok){          
                document.getElementById('v'+did).remove();
            }
        }catch{};       

    });
});

document.querySelectorAll('a[id^="e_"]').forEach(item => {
    item.addEventListener('click', async function(e) {
        e.preventDefault();        
        eid=item.getAttribute('id').substring(2);   
        const f=document.getElementById('f_' + eid ); 
        const fg=f.elements['a_' + eid];
        fg.value = fg.value.trim();
        if(fg.value.length < 5){
            document.getElementById("stat").innerHTML='At least 5 characters to save'; 
            return
        }
        
        try {
            const response = await fetch("/medit?id="+eid , {method:"Post", body: new FormData(f)});
            if (response.ok){
              document.getElementById("stat").innerHTML='Saving OK';                  
            }
        }catch{};             
  
    });
});