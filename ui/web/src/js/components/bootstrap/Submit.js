import React from "react";

class Submit extends React.Component {
	render() {
		return (
			// jshint ignore:start
			<input className="btn btn-default" type="submit"
				value={this.props.text} />
			// jshint ignore:end
		);
	}
}

module.exports = Submit;
