const apiurl = "http://localhost:5000/todos"

function markcompleted(id) {
    let confirmMC = confirm("Mark this todo as complete? You cannot edit once you mark a todo complete");
    if (confirmMC == true) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
                setTimeout(function () { window.location.reload(); }, 200);
            }
            else {
                console.log(this.responseText)
            }
        };
        xhttp.open("PUT", apiurl + "/markcompleted/" + id, true);
        xhttp.setRequestHeader("Content-type", "application/json");
        xhttp.send();
    }
}

// Function to delete a note
function deleteTodo(id) {
    let confirmDel = confirm("Delete this todo?");
    if (confirmDel == true) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
                setTimeout(function () { window.location.reload(); }, 200);
            }
            else {
                console.log(this.responseText)
            }
        };
        xhttp.open("DELETE", apiurl + "/" + id, true);
        xhttp.setRequestHeader("Content-type", "application/json");
        xhttp.send();
    }
}

function getTodos() {
    const url = "http://localhost:5000/todos"
    fetch(url)
        .then((response) => {
            return response.json()
        })
        .then((data) => {
            // Work with JSON data here
            if (data == null) {
                todosObj = [];
            } else {
                todosObj = data
            }
            console.log("todoObj:", todosObj)
            let html = "";
            if(todosObj["activetodos"])
            todosObj["activetodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <h4 class="todo-title"> ${element.title} </h4>
                    <p class="note-text"> ${element.text} </p>
                    <button id="${element.id}" type="button" class="todo-btn btn-primary" data-bs-toggle="modal" data-bs-target="#exampleModal"
                    data-bs-whatever="${element.title}" title-data="${element.title}" text-data="${element.text}" id-data=${element.id}>Edit Todo</button>
                    <button id="${element.id}"onclick="markcompleted(this.id)" class="todo-btn btn-success">Mark as completed</button>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="todo-btn btn-danger">Delete Todo</button>
                </div>
                    `;
            });
            let activeTodosElm = document.getElementById("activeTodos");
            if (todosObj['activetodos'].length != 0) {
                activeTodosElm.innerHTML = html;
            } else {
                activeTodosElm.innerHTML = `No Active Todos found!`;
            }

            html = "";
            todosObj["completedtodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <h4 class="todo-title"> ${element.title} </h4>
                    <p class="note-text"> ${element.text}</p>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="todo-btn btn-danger">Delete Todo</button>
                </div>
                    `;
            });
            let completedTodosElm = document.getElementById("completedTodos");
            if (todosObj['completedtodos'].length != 0) {
                completedTodosElm.innerHTML = html;
            } else {
                completedTodosElm.innerHTML = `No Completed Todos found!`;
            }

        })
        .catch((err) => {
            console.log(err)
        })

}

getTodos()


function IsEmpty(value) {
    return value ? value.trim().length == 0 : true;
}

var exampleModal = document.getElementById('exampleModal')
exampleModal.addEventListener('show.bs.modal', function (event) {
  // Button that triggered the modal
  var button = event.relatedTarget

  var id = button.getAttribute('id')
  console.log(id)

  if(id === "createbtn") {
    document.getElementById("exampleModalLabel").innerHTML = `Create Todo`
    document.getElementById("modalSubmit").innerHTML = `Create`
    console.log(true);
  } else {
    document.getElementById("exampleModalLabel").innerHTML = `Update Todo`
    document.getElementById("modalSubmit").innerHTML = `Update`
  }

  var modalTitle = exampleModal.querySelector('.modal-title')
  var modalBodyInputTitle = exampleModal.querySelector('.modal-body input')
  var modalBodyInputText = exampleModal.querySelector('.modal-body textarea')
  modalBodyInputTitle.value = button.getAttribute("title-data")
  modalBodyInputText.value = button.getAttribute("text-data")


  var submitBtn = document.getElementById("modalSubmit")
  submitBtn.addEventListener("click", function(ev) {
    if (IsEmpty(modalBodyInputTitle.value) || IsEmpty(modalBodyInputText.value)){
        alert("Fields cannot be empty")
    } else {

        var todojson = new Object();
        todojson.title = modalBodyInputTitle.value;
        todojson.text  = modalBodyInputText.value;
        todojson.completed = false;
        var jsonString= JSON.stringify(todojson);
        console.log(jsonString);

        if(id === "createbtn") {
            console.log("Inside Create"+ button.getAttribute('id'));
            var createrequest = new XMLHttpRequest();
            createrequest.onreadystatechange = function() {
                if (this.readyState == 4 && this.status == 200) {
                    console.log(this.responseText)
                    setTimeout(function () { window.location.reload(); }, 200);
                }
            };
            createrequest.open("POST", apiurl, true);
            createrequest.setRequestHeader("Content-type", "application/json");
            createrequest.send(jsonString);
        }
        else {
            console.log(button.getAttribute('id-data'));
            var updaterequest = new XMLHttpRequest();
            updaterequest.onreadystatechange = function() {
                if (this.readyState == 4 && this.status == 200) {
                    console.log(this.responseText)
                    setTimeout(function () { window.location.reload(); }, 200);
                }
            };
            updaterequest.open("PUT", apiurl + "/" + button.getAttribute('id'), true);
            updaterequest.setRequestHeader("Content-type", "application/json");
            updaterequest.send(jsonString);
        }
    }
  })

})
