document.getElementById('studentfrm').addEventListener('submit',
     async (e) => {	
        e.preventDefault(); 
        const f=e.target;
        f.elements['parentid'].value = document.getElementById('main').dataset.parentid; 
        f.elements['parentfullname'].value = document.getElementById('main').dataset.parentFullname;       
        const formData = new FormData(f); 
        try { 
            response = await fetch('/add/students', { method: 'POST', body: formData}); 
            data = await response.text();

            if (response.ok) { 
               const box1= document.getElementById('box1');
               const child = document.createElement('div');
               child.innerText = f.elements['firstname'].value  + ' ' + f.elements['lastname'].value;
               box1.appendChild(child);

                f.reset();
            }else{
                alert(data)
            }
         } catch (error) { };
        });