/**
 * Called the input field changes. Called for every key
 * stroke.
 */
function PostInput() {
	const p = document.getElementById('postinput')

	// https://stackoverflow.com/questions/2803880
	p.style.height = '';
	p.style.height = p.scrollHeight + 'px';
}
