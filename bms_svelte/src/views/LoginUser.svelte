<script>
    import { onMount } from "svelte";
    import Dialer, { AppendTree, ArgsConstructor, ElementConstructor } from "../utils";
    import { router } from "tinro";

let authForm, submitLogin, 
	inputs = inputObj => {
	return ElementConstructor({
	el: "input",
	placeholder: inputObj.label,
	type: inputObj.type
	}, {"margin-bottom": "24px"});
}

onMount(() => {
	let username = inputs({ label: "Name or Work Id", type: "text" }), password = inputs({ label: "Password", type: "password" });

	submitLogin.onclick = () => {
		const usernameVal = username.children[1].value.trim(), passwordVal = password.children[1].value.trim();

		if (usernameVal.length < 1 || passwordVal.length < 1)
			ElementConstructor({ el: "dialog", dialogType: "error", dialogError: "Fields cannot be empty" });
		else
			Dialer(Login, ArgsConstructor(["", [[usernameVal], [passwordVal]]])).then(res => {
				if (res.Err != undefined)
					ElementConstructor({ el: "dialog", dialogType: "error", dialogError: res.Err });
				else if (document.URL.includes("/login"))
					router.goto("/");
				else window.location.assign("/");
			}); 
	}

	AppendTree([authForm], [username, password]);
});
</script>

<main>
	<div class="bms_auth_div">
		<p style="margin: 0 0 36px 0;"><span class="bms_header"><b>beazy</b></span> <small>login</small></p>
		<div class="bms_auth_form" bind:this={authForm}></div>
		<button type="submit" bind:this={submitLogin}>Login</button>
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
