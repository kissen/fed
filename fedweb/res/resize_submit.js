/**
 * Called the input field changes. Called for every key
 * stroke.
 */
function SubmitInput() {
	const submit = document.getElementById('submit')

	// https://stackoverflow.com/questions/2803880
	submit.style.height = '';
	submit.style.height = submit.scrollHeight + 'px';
}
