// const apiurl = "http://localhost:8000/todos"

const apiurl = "https://praveenb-todoapp.herokuapp.com/todos"

// Function to mark a todo as complete
function markcompleted(id) {
    let confirmMC = confirm("Mark this todo as complete? You cannot edit once you mark a todo complete");
    if (confirmMC == true) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
                // window.location.reload();
                getTodos();
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

// Function to delete a todo
function deleteTodo(id) {
    let confirmDel = confirm("Delete this todo?");
    if (confirmDel == true) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
                // window.location.reload();
                getTodos();
            }
        };
        xhttp.open("DELETE", apiurl + "/" + id, true);
        xhttp.setRequestHeader("Content-type", "application/json");
        xhttp.send();
    }
}

// Function to get all todos
function getTodos() {
    fetch(apiurl)
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

            // console.log(todosObj)

            let activeTodosElm = document.getElementById("activeTodos");
            let completedTodosElm = document.getElementById("completedTodos");

            activeTodosElm.innerHTML = ``;
            completedTodosElm.innerHTML = ``;
            

            let html = "";
            if(todosObj["activetodos"])
            todosObj["activetodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <p class="todo-date text-secondary mx-1"><b>Created on:</b> ${element.date}</p>
                    <h3 class="todo-title m-1"> ${element.title} </h3>
                    <p class="todo-text lead mx-1 mt-1 mb-3"> ${element.text} </p>
                    <button id="${element.id}" type="button" class="btn todo-btn btn-primary m-1" data-bs-toggle="modal" data-bs-target="#exampleModal"
                    data-bs-whatever="${element.title}" title-data="${element.title}" text-data="${element.text}" id-data=${element.id}>Edit Todo</button>
                    <button id="${element.id}"onclick="markcompleted(this.id)" class="btn todo-btn btn-success m-1">Mark as completed</button>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="btn todo-btn btn-danger m-1">Delete Todo</button>
                </div>
                    `;
            });
            
            if (todosObj['activetodos'] != null && todosObj['activetodos'].length != 0) {
                activeTodosElm.innerHTML = html;
            } else {
                activeTodosElm.innerHTML = `<p class="lead px-1">No todos found!</p>`;
            }

            html = "";
            if(todosObj["completedtodos"])
            todosObj["completedtodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <p class="todo-date text-secondary"><b>Created on</b> ${element.date}</p>
                    <h4 class="todo-title"> ${element.title} </h4>
                    <p class="todo-text lead"> ${element.text}</p>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="btn todo-btn btn-danger m-1">Delete Todo</button>
                </div>
                    `;
            });
            
            if (todosObj['completedtodos'] != null && todosObj['completedtodos'].length != 0) {
                completedTodosElm.innerHTML = html;
            } else {
                completedTodosElm.innerHTML = `<p class="lead px-1">No todos found!</p>`;
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
  event.stopPropagation();
  // Button that triggered the modal
  var button = event.relatedTarget

  var id = button.getAttribute('id')
  console.log(id);

  if(id === "createbtn") {
    document.getElementById("exampleModalLabel").innerHTML = `Create Todo`
    document.getElementById("modalSubmit").innerHTML = `Create`
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

        if(id === "createbtn") {
            console.log(button.getAttribute('id'));
            var createrequest = new XMLHttpRequest();
            createrequest.onreadystatechange = function() {
                if (this.readyState == 4 && this.status == 200) {
                    console.log(this.responseText);
                    window.location = window.location
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
                    console.log(this.responseText);

                }
            };
            updaterequest.open("PUT", apiurl + "/" + button.getAttribute('id'), true);
            updaterequest.setRequestHeader("Content-type", "application/json");
            updaterequest.send(jsonString);
        }
    }
  })

})

var searchInputField = document.getElementById("searchinputfield")
searchInputField.addEventListener("keyup", function(e) {
    if (e.key === 'Enter') {
        e.preventDefault();
        document.getElementById("searchbtn").click();
    }
});

var searchBtn = document.getElementById("searchbtn")
searchBtn.addEventListener("click", function(ev) {

    ev.preventDefault();

    if (IsEmpty(searchInputField.value)){
        alert("Search query should not be empty!")
    } else {
        fetch(apiurl + "/search?query=" + searchInputField.value)
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

            let newTododElement = document.getElementById("newTodo");
            newTododElement.style.display = "none";

            let searchResultElm = document.getElementById("searchResultDiv");
            let activeTodosElm = document.getElementById("activeTodos");
            let completedTodosElm = document.getElementById("completedTodos");

            let searchcount = 0;
            if(todosObj["activetodos"] && todosObj['activetodos'] != null) {
                searchcount += todosObj['activetodos'].length;
            }
            if(todosObj["completedtodos"] && todosObj['completedtodos'] != null) {
                searchcount += todosObj['completedtodos'].length;
            }
            searchResultElm.innerHTML = `<p class="lead"><b>${searchcount} results found.</b> <a style="color:#0790D1"; href="/"> Return to main menu</a></p>`;


            activeTodosElm.innerHTML = ``;
            completedTodosElm.innerHTML = ``;
            

            let html = "";
            if(todosObj["activetodos"])
            todosObj["activetodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <p class="todo-date text-secondary mx-1"><b>Created on:</b> ${element.date}</p>
                    <h3 class="todo-title m-1"> ${element.title} </h3>
                    <p class="todo-text lead mx-1 mt-1 mb-3"> ${element.text} </p>
                    <button id="${element.id}" type="button" class="btn todo-btn btn-primary m-1" data-bs-toggle="modal" data-bs-target="#exampleModal"
                    data-bs-whatever="${element.title}" title-data="${element.title}" text-data="${element.text}" id-data=${element.id}>Edit Todo</button>
                    <button id="${element.id}"onclick="markcompleted(this.id)" class="btn todo-btn btn-success m-1">Mark as completed</button>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="btn todo-btn btn-danger m-1">Delete Todo</button>
                </div>
                    `;
            });
            
            if (todosObj['activetodos'] != null && todosObj['activetodos'].length != 0) {
                activeTodosElm.innerHTML = html;
            } else {
                activeTodosElm.innerHTML = `<p class="lead px-1">No todos found!</p>`;
            }

            html = "";
            if(todosObj["completedtodos"])
            todosObj["completedtodos"].reverse().forEach(function (element, index) {
                html += `
                <div class="todo">
                    <p class="todo-date text-secondary"><b>Created on</b> ${element.date}</p>
                    <h4 class="todo-title"> ${element.title} </h4>
                    <p class="todo-text lead"> ${element.text}</p>
                    <button id="${element.id}"onclick="deleteTodo(this.id)" class="btn todo-btn btn-danger m-1">Delete Todo</button>
                </div>
                    `;
            });
            
            if (todosObj['completedtodos'] != null && todosObj['completedtodos'].length != 0) {
                completedTodosElm.innerHTML = html;
            } else {
                completedTodosElm.innerHTML = `<p class="lead px-1">No todos found!</p>`;
            }

        })
        .catch((err) => {
            console.log(err)
        })
    }

})
