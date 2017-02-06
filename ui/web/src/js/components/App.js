import React from "react";
import Feeds from "./Feeds";
import ContextMenu from "./ContextMenu";

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
    return (
      // jshint ignore:start
      <div className="row">
        <div className="col-md-3">
          <Feeds data={this.props.feeds} onSelectFeed={this.selectFeed}
                 onFeedContextMenu={this.showFeedContextMenu} />
        </div>

        <div className="col-md-9">
          <div className="row">
          </div>
          <div className="row">
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
  onSelectFeed: React.PropTypes.func.isRequired
};

module.exports = App;
