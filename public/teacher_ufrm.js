     document.querySelectorAll("#ffrm").forEach(frm => {
        frm.addEventListener('submit', async (e) => {
            e.preventDefault();        
            const fdata=new FormData(e.target);        
            const response=await fetch('/ut', {method: 'Post', body: fdata });
             const data=await response.text();
                if (data === "OK"){
                    alert("OK Saving")
                }else{
                    alert("Error Saving")
                }
            });
        })
