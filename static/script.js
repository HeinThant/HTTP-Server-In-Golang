const API_URL = "/items";

async function fetchItems() {
    const response = await fetch(API_URL);
    const items = await response.json();
    const list = document.getElementById("itemList");
    list.innerHTML = "";
    items.forEach(item => {
        let li = document.createElement("li");
        li.textContent = `${item.id}: ${item.name}`;
        
        let deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.onclick = () => deleteItem(item.id);

        let updateBtn = document.createElement("button");
        updateBtn.textContent = "Update";
        updateBtn.onclick = () => updateItem(item.id);

        li.appendChild(updateBtn);
        li.appendChild(deleteBtn);
        list.appendChild(li);
    });
}

async function addItem() {
    const name = document.getElementById("itemName").value;
    if (!name) return;
    
    await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
    });
    document.getElementById("itemName").value = "";
    fetchItems();
}

async function deleteItem(id) {
    await fetch(`${API_URL}?id=${id}`, { method: "DELETE" });
    fetchItems();
}

async function updateItem(id) {
    const newName = prompt("Enter new name:");
    if (!newName) return;

    await fetch(`${API_URL}?id=${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: newName }),
    });
    fetchItems();
}

async function resetServer() {
    await fetch("/reset", { method: "POST" });
    fetchItems(); // Refresh item list
}

fetchItems();
