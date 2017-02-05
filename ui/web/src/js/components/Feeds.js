import React from "react";

class Feed extends React.Component {
	constructor(props) {
		super(props);
		this.selectFeed = this.selectFeed.bind(this);
	}

	className() {
		var klass = ['feed'];
		if(this.props.data.Selected) {
			klass.push('selected');
		}
		return klass.join(' ');
	}

	selectFeed() {
		this.props.onSelectFeed(this.props.data);
	}

	render() {
		return (
			// jshint ignore:start
			<div onClick={this.selectFeed} className={this.className()}>
				<span>{this.props.data.Title}</span>
			</div>
			// jshint ignore:end
		);
	}
}

Feed.propTypes = {
	feed: React.PropTypes.object
};

class Feeds extends React.Component {
	render() {
		var feeds = this.props.data.map((f, i) => {
			return (
				// jshint ignore:start
				<li key={i.toString()}>
					<Feed onSelectFeed={this.props.onSelectFeed} data={f} />
				</li>
				// jshint ignore:end
			);
		});
		// jshint ignore:start
		return <ul>{feeds}</ul>;
		// jshint ignore:end
	}
}

Feeds.propTypes = {
	data: React.PropTypes.array,
	onSelectFeed: React.PropTypes.func
};

Feeds.Feed = Feed;
module.exports = Feeds;
