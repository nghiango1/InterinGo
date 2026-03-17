<script lang="ts">
	import { page } from '$app/state';

	const status = page.status;
	const message = page.error?.message;
</script>

<div class="container">
	<h1 class="header">{status}</h1>
	<div class="icon">
		<i class="fas fa-exclamation-triangle"></i>
	</div>

	{#if message}
		<p>Oops! We try to find the page you are looking for, but got error:</p>
		<pre class="text-center">{message}</pre>
	{/if}
	<a href="/">Go back to the Home Page</a>
</div>

<style>
	h1 {
		color: #ff4d4d;
		text-align: center;
		margin-top: 20px;
		font-size: 40px;
	}

	p {
		text-align: center;
		margin: 10px;
	}

	a {
		display: block;
		text-align: center;
		margin-top: 20px;
		color: #007acc;
		text-decoration: none;
	}

	a:hover {
		text-decoration: underline;
	}

	.container {
		max-width: 600px;
		margin: 0 auto;
	}

	.icon {
		font-size: 100px;
		color: #ff4d4d; /* Match icon color with heading */
		text-align: center;
		margin-bottom: 20px;
	}

	@media (max-width: 768px) {
		.container {
			padding: 0 0.5rem;
		}

		.header {
			flex-direction: column;
			align-items: center;
			text-align: center;
			padding: 1.5rem;
			gap: 1.5rem;
		}
	}
</style>
