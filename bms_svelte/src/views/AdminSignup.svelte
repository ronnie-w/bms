<script>
    import { onMount } from "svelte";
    import Dialer, { AppendTree, ArgsConstructor, ElementConstructor } from "../utils";

export let signUpMessage;

let authForm, submitCredentials,
	inputs = inputObj => {
	return ElementConstructor({
		el: "input",
		placeholder: inputObj.label,
		type: inputObj.type,
		readonly: inputObj.readonly,
		value: inputObj.value,
	}, { "margin-bottom": "24px" });
}

onMount(() => {
	let username = inputs({ 
		label: "Name", 
		readonly: "readonly",
		value: "admin"
	}), email = inputs({ label: "Email" }), password = inputs({ label: "Password", type: "password" });

	submitCredentials.onclick = () => {
		if (password.children[1].value.length < 4)
			ElementConstructor({ el: "dialog", dialogType: "error", dialogError: "Password length must be four characters or more" })
		if (!/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/.test(email.children[1].value.trim()))
			ElementConstructor({ el: "dialog", dialogType: "error", dialogError: "The email you entered is invalid" })
		else {
			Dialer(Signup, ArgsConstructor(["admin", [
				[username.children[1].value.trim()], [email.children[1].value.trim()], [password.children[1].value.trim()], ["0"], ["0"]
			]])).then(res => {
				if (res.Err != undefined)
					ElementConstructor({ el: "dialog", dialogType: "error", dialogError: res.Err })
				else window.location.assign("/login");
			});
		}
	}

	AppendTree([authForm], [username, email, password]);

});
</script>

<main>
	<div class="bms_auth_div">
		<p style="margin: 0 0 6px 0;"><span class="bms_header"><b>beazy</b></span> <small>create admin account</small></p>
		<small style="color: grey;">{signUpMessage}</small>
		<div class="bms_auth_form" style="margin-top: 12px;" bind:this={authForm}></div>
		<p><a href="/reset_password">Forgot password?</a></p>
		<button type="submit" bind:this={submitCredentials}>Create admin</button>
	</div>
</main>

<style>
	.bms_auth_div {
		box-shadow: rgba(60, 64, 67, 0.3) 0px 1px 2px 0px, rgba(60, 64, 67, 0.15) 0px 2px 6px 2px;
		width: fit-content;
		padding: 36px;
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		border: 6px;
	}

	.bms_header {
		font-family: "Comfortaa", cursive;
		cursor: pointer;
		user-select: none;
	}
</style>
