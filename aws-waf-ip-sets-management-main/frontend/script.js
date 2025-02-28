document.getElementById('createIpSetForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const name = document.getElementById('createName').value;
    const addresses = document.getElementById('createAddresses').value.split(',').map(addr => addr.trim());
    const description = document.getElementById('createDescription').value;

    if (confirm("Are you sure you want to create this IP set?")) {
        const userInput = prompt("Type 'Confirm' to proceed with creating the IP set:");
        if (userInput === "Confirm") {
            try {
                const response = await fetch('/api/create-ip-set', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name, addresses, description })
                });

                const data = await response.json();
                if (response.ok) {
                    console.log(data);
                    fetchIpSets(); // Refresh the list after creation
                } else {
                    console.error('Error creating IP set:', data.error);
                    alert('Error creating IP set: ' + data.error);
                }
            } catch (error) {
                console.error('Fetch error:', error);
                alert('Fetch error: ' + error);
            }
        } else {
            alert("Operation canceled.");
        }
    }
});

document.getElementById('deleteIpSetForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const id = document.getElementById('deleteId').value;
    const name = document.getElementById('deleteName').value;
    const lockToken = document.getElementById('deleteLockToken').value;

    if (confirm("Are you sure you want to delete this IP set?")) {
        const userInput = prompt("Type 'Confirm' to proceed with deleting the IP set:");
        if (userInput === "Confirm") {
            try {
                const response = await fetch('/api/delete-ip-set', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ id, name, lockToken })
                });

                const data = await response.json();
                if (response.ok) {
                    console.log(data);
                    fetchIpSets(); // Refresh the list after deletion
                } else {
                    console.error('Error deleting IP set:', data.error);
                    alert('Error deleting IP set: ' + data.error);
                }
            } catch (error) {
                console.error('Fetch error:', error);
                alert('Fetch error: ' + error);
            }
        } else {
            alert("Operation canceled.");
        }
    }
});

document.getElementById('addIpAddressForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (confirm("Are you sure you want to add these IP addresses?")) {
        const id = document.getElementById('addId').value;
        const name = document.getElementById('addName').value;
        const addresses = document.getElementById('addAddresses').value.split(',').map(addr => addr.trim());
        const lockToken = document.getElementById('addLockToken').value;

        try {
            const response = await fetch('/api/add-ip-address', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ id, name, addresses, lockToken })
            });

            const data = await response.json();
            if (response.ok) {
                console.log(data);
                fetchIpSets(); // Refresh the list after update
            } else {
                console.error('Error adding IP address:', data.error);
                alert('Error adding IP address: ' + data.error);
            }
        } catch (error) {
            console.error('Fetch error:', error);
            alert('Fetch error: ' + error);
        }
    }
});

document.getElementById('removeIpAddressForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (confirm("Are you sure you want to remove these IP addresses?")) {
        const id = document.getElementById('removeId').value;
        const name = document.getElementById('removeName').value;
        const addresses = Array.from(document.getElementById('removeAddresses').selectedOptions).map(option => option.value);
        const lockToken = document.getElementById('removeLockToken').value;

        try {
            const response = await fetch('/api/remove-ip-address', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ id, name, addresses, lockToken })
            });

            const data = await response.json();
            if (response.ok) {
                console.log(data);
                fetchIpSets(); // Refresh the list after update
            } else {
                console.error('Error removing IP address:', data.error);
                alert('Error removing IP address: ' + data.error);
            }
        } catch (error) {
            console.error('Fetch error:', error);
            alert('Fetch error: ' + error);
        }
    }
});

async function fetchIpSets() {
    try {
        const response = await fetch('/api/list-ip-sets');
        const data = await response.json();
        const ipSetsDiv = document.getElementById('ipSets');
        const addSelect = document.getElementById('addSelect');
        const removeSelect = document.getElementById('removeSelect');
        const deleteSelect = document.getElementById('deleteSelect');
        
        ipSetsDiv.innerHTML = ''; // Clear the div first
        addSelect.innerHTML = '<option value="" disabled selected>Select an IP Set</option>'; // Clear the add select options first
        removeSelect.innerHTML = '<option value="" disabled selected>Select an IP Set</option>'; // Clear the remove select options first
        deleteSelect.innerHTML = '<option value="" disabled selected>Select an IP Set</option>'; // Clear the delete select options first

        data.forEach(ipSet => {
            const ipSetDiv = document.createElement('div');
            ipSetDiv.className = 'card mb-3 col-md-12';
            ipSetDiv.innerHTML = `
                <div class="card-body">
                    <h5 class="card-title">${ipSet.Name}</h5>
                    <p class="card-text"><strong>ID:</strong> ${ipSet.Id}</p>
                    <div class="ip-list">
                        ${ipSet.Addresses.map(addr => `<p>${addr}</p>`).join('')}
                    </div>
                    <p class="card-text"><strong>Total IPs:</strong> ${ipSet.AddressCount}</p>
                    <p class="card-text"><strong>LockToken:</strong> ${ipSet.LockToken}</p>
                </div>
            `;
            ipSetsDiv.appendChild(ipSetDiv);

            const addOption = document.createElement('option');
            addOption.value = ipSet.Id;
            addOption.textContent = ipSet.Name;
            addOption.dataset.lockToken = ipSet.LockToken;
            addSelect.appendChild(addOption);

            const removeOption = document.createElement('option');
            removeOption.value = ipSet.Id;
            removeOption.textContent = ipSet.Name;
            removeOption.dataset.lockToken = ipSet.LockToken;
            removeOption.dataset.addresses = JSON.stringify(ipSet.Addresses);
            removeSelect.appendChild(removeOption);

            const deleteOption = document.createElement('option');
            deleteOption.value = ipSet.Id;
            deleteOption.textContent = ipSet.Name;
            deleteOption.dataset.lockToken = ipSet.LockToken;
            deleteSelect.appendChild(deleteOption);
        });

        addSelect.addEventListener('change', (e) => {
            const selectedOption = addSelect.options[addSelect.selectedIndex];
            const selectedId = selectedOption.value;
            const selectedName = selectedOption.textContent;
            const selectedLockToken = selectedOption.dataset.lockToken;

            document.getElementById('addId').value = selectedId;
            document.getElementById('addName').value = selectedName;
            document.getElementById('addLockToken').value = selectedLockToken;
        });

        removeSelect.addEventListener('change', (e) => {
            const selectedOption = removeSelect.options[removeSelect.selectedIndex];
            const selectedId = selectedOption.value;
            const selectedName = selectedOption.textContent;
            const selectedLockToken = selectedOption.dataset.lockToken;
            const addresses = JSON.parse(selectedOption.dataset.addresses);

            document.getElementById('removeId').value = selectedId;
            document.getElementById('removeName').value = selectedName;
            document.getElementById('removeLockToken').value = selectedLockToken;

            const removeAddressesSelect = document.getElementById('removeAddresses');
            removeAddressesSelect.innerHTML = '<option value="" disabled>Select IP Addresses</option>';
            addresses.forEach(address => {
                const option = document.createElement('option');
                option.value = address;
                option.textContent = address;
                removeAddressesSelect.appendChild(option);
            });
        });

        deleteSelect.addEventListener('change', (e) => {
            const selectedOption = deleteSelect.options[deleteSelect.selectedIndex];
            const selectedId = selectedOption.value;
            const selectedName = selectedOption.textContent;
            const selectedLockToken = selectedOption.dataset.lockToken;

            document.getElementById('deleteId').value = selectedId;
            document.getElementById('deleteName').value = selectedName;
            document.getElementById('deleteLockToken').value = selectedLockToken;
        });
    } catch (error) {
        console.error('Fetch error:', error);
        alert('Fetch error: ' + error);
    }
}

fetchIpSets();
