import React from "react";
import Feeds from "./Feeds";
import Items from "./Items";
import ContextMenu from "./ContextMenu";
import Subscribe from "./Subscribe";

class App extends React.Component {
	constructor(props) {
		super(props);

		this.selectFeed = (data) => {
			this.props.onSelectFeed(data);
		};

		this.mouseX = this.mouseY = 0;

		this.saveMouseCoordinates = this.saveMouseCoordinates.bind(this);

		this.state = {
			showContextMenu: false,
			contextMenuOptions: {},
			contextMenuCallback: function() {}
		};
	}

	saveMouseCoordinates(e) {
		this.mouseY = e.clientY;
		this.mouseX = e.clientX;
	}

	showFeedContextMenu(data) {
		this.setState({
			showContextMenu: true,
			contextMenuOptions: {
				edit: "Edit Feed"
			},
			contextMenuCallback: (key) => {
				this.hideContextMenu();
			}
		});
	}

	hideContextMenu() {
		this.setState({
			showContextMenu: false,
			contextMenuOptions: {},
			contextMenuCallback: function() {}
		});
	}

	componentWillMount() {
		window.document.addEventListener('mousemove', this.saveMouseCoordinates);
	}

	componentWillUnmount() {
		window.document.removeEventListener('mousemove', this.saveMouseCoordinates);
	}

	render() {
		let items = [];

		for(let feed of this.props.feeds) {
			if(feed.Selected) {
				items = feed.Items;
				console.log(items);
				break;
			}
		}

		return (
			// jshint ignore:start
			<div className="row">
			<div className="col-md-3">
			<Feeds data={this.props.feeds} onSelectFeed={this.selectFeed}
			onFeedContextMenu={this.showFeedContextMenu} />
			</div>

			<div className="col-md-9">
			<div className="row">
			<Subscribe onSubmit={this.props.onSubscribeToFeed} />
			</div>
			<div className="row">
			<Items data={items} onSelectItem={this.props.onSelectItem} />
			</div>
			</div>

			<ContextMenu top={this.mouseY} left={this.mouseX}
			show={this.state.showContextMenu} options={this.state.contextMenuOptions}
			onSelectOption={this.state.contextMenuCallback}/>

			</div>
			// jshint ignore:end
		);
	}
}

App.propTypes = {
	feeds: React.PropTypes.array.isRequired,
	onSelectFeed: React.PropTypes.func.isRequired,
	onSubscribeToFeed: React.PropTypes.func.isRequired,
	onSelectItem: React.PropTypes.func.isRequired
};

module.exports = App;
