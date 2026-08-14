document.getElementById('uc').addEventListener('click', function(e){
    e.preventDefault();
    const p=document.getElementById('uf');
    p.remove();
 });

 document.getElementById('frmudg').addEventListener('submit',
    async (e) => {	
       e.preventDefault(); 
       const f=e.target;      
       f.elements['grade'].value =   f.elements['grade'].value.trim();     
       const formData = new FormData(f); 
       try { 
           response = await fetch('/udg', { method: 'POST', body: formData}); 
           data = await response.text();

           if (response.ok) { 
            alert('Update grade Successfully');
             const pp=document.getElementById('uf');
             pp.remove();
           }else{
               alert(data);
           }
        } catch (error) { };
       });