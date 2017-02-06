import React from "react";
import $ from "jquery";

class Modal extends React.Component {
  constructor(props) {
    super(props);
    this.state = { show: props.show };
    this.handleManualHide = this.handleManualHide.bind(this);
  }

  showHide() {
    if(this.element) {
      $(this.element).modal(this.state.show ? 'show': 'hide');
    }
  }

  componentDidMount() {
    this.showHide();
    $(this.element).on('hidden.bs.modal', this.handleManualHide);
  }

  componentWillUnmount() {
    $(this.element).off('hidden.bs.modal', this.handleManualHide);
  }

  handleManualHide() {
    this.setState({show: false});
  }

  componentWillReceiveProps(newProps) {
    if(newProps.show != this.props.show) {
      this.setState({show: newProps.show}, () => { this.showHide(); });
    }
  }

  render() {
    return (
    // jshint ignore:start
      <div ref={(el) => { this.element = el; }} className="modal fade" tabIndex="-1" role="dialog">
        <div className="modal-dialog" role="document">
          <div className="modal-content">
            <div className="modal-header">
              <button onClick={this.close} type="button" className="close" data-dismiss="modal"
                aria-label="Close"><span aria-hidden="true">&times;</span></button>
              <h4 className="modal-title">{this.props.title}</h4>
            </div>
            <div className="modal-body">{this.props.body}</div>
            <div className="modal-footer">{this.props.footer}</div>
          </div>
        </div>
      </div>
    // jshint ignore:end
    );
  }
}

module.exports = Modal;
