<script>
    import { onMount } from "svelte";
    import Dialer, { AppendTree, ArgsConstructor, ElementConstructor } from "../utils";

let authForm, submitCredentials;

onMount(() => {
	let email = inputs({ label: "Email" }), 
		password = inputs({ label: "New password", type: "password" });

	submitCredentials.onclick = () => {
		if (password.children[1].value.length < 4)
			ElementConstructor({ el: "dialog", dialogType: "error", dialogError: "Password length must be four characters or more" })
		else if (!/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/.test(email.children[1].value.trim()))
			ElementConstructor({ el: "dialog", dialogType: "error", dialogError: "The email you entered is invalid" })
		else Dialer(ResetPassword, ArgsConstructor(["admin", [[email.children[1].value.trim()], [password.children[1].value.trim()]]]))
				.then(res => {
						console.log(res);
						if (res.Err)
							ElementConstructor({el: "dialog", dialogType: "error", dialogError: res.Err});
						else window.location.assign("/login");
				});
	}

	AppendTree([authForm], [email, password]);

});

function inputs(inputObj) {
	return ElementConstructor({
		el: "input",
		placeholder: inputObj.label,
		type: inputObj.type
	}, { "margin-bottom": "24px" });
}
</script>

<main>
	<div class="bms_auth_div">
		<p style="margin: 0 0 36px 0;"><span class="bms_header"><b>beazy</b></span> <small>reset password</small></p>
		<div class="bms_auth_form" bind:this={authForm}></div>
		<button type="submit" bind:this={submitCredentials}>Submit</button>
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
