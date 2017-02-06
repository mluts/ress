import React from "react";

class Option extends React.Component {
	constructor(props) {
		super(props);
		this.selectOption = this.selectOption.bind(this);
	}

	selectOption() {
		this.props.callback(this.props.optionKey);
	}

	render() {
		return (
			// jshint ignore:start
			<li onClick={this.selectOption}>
				{this.props.title}
			</li>
			// jshint ignore:end
		);
	}
}

Option.propTypes = {
	optionKey: React.PropTypes.string,
	title: React.PropTypes.string,
	callback: React.PropTypes.func,
};

class ContextMenu extends React.Component {
	constructor(props) {
		super(props);
		this.state = { show: props.show };
		this.selectOption = this.selectOption.bind(this);
	}

	style() {
		return {
			position: "fixed",
			top: this.props.top,
			left: this.props.left
		};
	}

	selectOption(data) {
		this.props.onSelectOption(data);
		this.setState({show: false});
	}

	componentWillReceiveProps(props) {
		this.setState({show: props.show});
	}

	render() {
		if(!this.state.show) {
			return null;
		}

		var options = [];

		for(let key in this.props.options) {
			options.push(
				// jshint ignore:start
				<Option key={key} optionKey={key} title={this.props.options[key]}
								callback={this.selectOption} />
				// jshint ignore:end
			);
		}

		return (
			// jshint ignore:start
			<div style={this.style()}>
				<ul>{options}</ul>
			</div>
			// jshint ignore:end
		);
	}
}

ContextMenu.propTypes = {
	top: React.PropTypes.number,
	left: React.PropTypes.number,
	options: React.PropTypes.object,
	onSelectOption: React.PropTypes.func,
	show: React.PropTypes.bool
};

module.exports = ContextMenu;
