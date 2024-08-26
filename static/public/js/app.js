const $ = (selector) => document.querySelector(selector);
const container = $('#users');
const API_ENDPOINT = '/api/v1/users';

const listUsers = async () => {
	const response = await fetch(API_ENDPOINT);
	const data = await response.json();
	const users = data.users.reverse();

	for (let index = 0; index < users.length; index++) {
		const child = document.createElement('li');
		child.className = 'list-group-item';
		child.innerText = users[index].name;

		container.appendChild(child);
	}
};

$('#add-user').addEventListener('click', async (e) => {
	e.preventDefault();
	const user = $('#user').value;

	if (!user) return;

	const form = new FormData();
	form.append('user', user);

	const response = await fetch(API_ENDPOINT, {
		method: 'POST',
		body: form,
	});

	const data = await response.json();

	const child = document.createElement('li');
	child.className = 'list-group-item';
	child.innerText = data.user.name;

	container.insertBefore(child, container.firstChild);

	$('#user').value = '';
});

$('#delete-user').addEventListener('click', async (e) => {
	e.preventDefault();

	const form = new FormData();
	form.append('user', "");

	const response = await fetch(API_ENDPOINT, {
		method: 'DELETE',
		body: form,
	});

	await response.json();
	container.innerHTML = "";

});

document.addEventListener('DOMContentLoaded', listUsers);
