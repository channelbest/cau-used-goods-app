<script>
	export default {
		onLaunch: function() {
			console.log('App Launch')
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style>
	page {
		background: #f6f8f5;
		color: #24302b;
		font-size: 28rpx;
	}

	button::after {
		border: 0;
	}

	.empty-state {
		padding: 96rpx 24rpx;
		color: #89938f;
		text-align: center;
	}
</style>
