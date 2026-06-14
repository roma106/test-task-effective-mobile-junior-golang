// Adding subscription
let createSubsButton = document.querySelector(".create-subs-button");

createSubsButton.addEventListener("click", CreateSubscription);

function CreateSubscription() {
    let data = {
        name: document.querySelector(".create-subs-name").value,
        price: Number(document.querySelector(".create-subs-price").value),
        user_id: document.querySelector(".create-subs-user").value,
        start_date: document.querySelector(".create-subs-start").value,
        end_date: document.querySelector(".create-subs-end").value,
    };

    fetch("/subs", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(res => {
            if (res.ok) {
                alert("Подписка добавлена");
                GetSubs();
            } else {
                res.text().then(data => alert(data));
            }
        });
}


// Geting all subs and render
const SubsContainer = document.querySelector(".subs-list");

GetSubs();

function GetSubs() {
    fetch("/subs")
        .then(res => res.json())
        .then(data => {
            SubsContainer.innerHTML = "";

            data.forEach(sub => {
                let newSub = document.createElement("div");
                newSub.classList.add("sub");

                newSub.innerHTML = `
                    <p class="sub-info">${sub.name}</p>
                    <p class="sub-info">${sub.price}</p>
                    <p class="sub-info">${sub.user_id}</p>
                    <p class="sub-info">${ParseDate(sub.start_date)}</p>
                    <p class="sub-info">${(sub.end_date == null) ? "не указано" : ParseDate(sub.end_date)}</p>
                    <p class="sub-button-edit" id="${sub.id}">Редактировать</p>
                    <p class="sub-button-edit-ready hidden" id="${sub.id}">Готово</p>
                    <p class="sub-button-delete" id="${sub.id}">Удалить</p>
                `;

                let editBtn = newSub.querySelector(".sub-button-edit");

                editBtn.addEventListener("click", () => {
                    MakeFieldsToEdit(editBtn);
                });

                let readyBtn = newSub.querySelector(".sub-button-edit-ready");

                readyBtn.addEventListener("click", () => {
                    EditSub(readyBtn)
                });
                newSub.querySelector(".sub-button-delete").addEventListener("click", ()=>{
                    DeleteSub(sub.id)
                })

                SubsContainer.appendChild(newSub);
            });
        });
}


function ParseDate(date) {
    let year = String(date).slice(0, 4);
    let month = String(date).slice(5, 7);
    // console.log("Please take me to your company, I am a cool guy")
    return month + "-" + year;
}

// Editing sub


function MakeFieldsToEdit(editBtn) {
    let subElement = editBtn.parentElement;

    editBtn.classList.add("hidden")
    let readyBtn = subElement.querySelector(".sub-button-edit-ready")
    readyBtn.classList.remove("hidden")


    subElement.querySelectorAll(".sub-info").forEach(el => {
        let input = document.createElement("input");

        input.classList.add("edit-sub-input");
        input.value = el.innerText;

        el.classList.add("hidden");
        el.parentElement.insertBefore(input, el);
    });
}

function EditSub(readyBtn){
    let inputs = readyBtn.parentElement.querySelectorAll(".edit-sub-input");
    let data = {
        id: Number(readyBtn.id),
        name: inputs[0].value,
        price: Number(inputs[1].value),
        user_id: inputs[2].value,
        start_date: inputs[3].value,
        end_date: (inputs[4].value=="не указано") ? "" : inputs[4].value,
    };

    fetch("/subs/"+String(data.id), {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(res => {
            if (res.ok) {
                alert("Подписка обновлена");
                GetSubs();
            } else {
                res.text().then(data => alert(data));
            }
        }); 
}

function DeleteSub(id){
    fetch("/subs/"+String(id), {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
    })
        .then(res => {
            if (res.ok) {
                alert("Подписка удалена");
                GetSubs();
            } else {
                res.text().then(data => alert(data));
            }
        });
}



// Sum prices

let openFilterButton = document.querySelector(".subs-filter")
let filterContainer = document.querySelector(".filter-container")
let closeFilterButton = document.querySelector(".filter-close-button")

openFilterButton.addEventListener("click", ()=>{
    filterContainer.classList.remove("hidden")
})
closeFilterButton.addEventListener("click", ()=>{
    filterContainer.classList.add("hidden")
})

let filterGoButton = document.querySelector(".filter-go")
filterGoButton.addEventListener("click", CountSum)
let answerContainer = document.querySelector(".filter-answer-container")

function CountSum(){
    let data = {
        name: document.querySelector(".filter-name").value,
        user_id: document.querySelector(".filter-user-id").value,
        start_date: document.querySelector(".filter-period-start").value,
        end_date: document.querySelector(".filter-period-end").value,
    }

    fetch("/subs/sum?name="+data.name+"&user_id="+data.user_id+"&period_start="+data.start_date+"&period_end="+data.end_date, {
        method: "GET"
    }).then(res =>{
        if (!res.ok){
            res.text().then(data => alert(data))
        }else{
            answerContainer.classList.remove("hidden")
            res.text().then(data => { document.querySelector(".filter-answer").innerHTML = data})
        }
    })
}