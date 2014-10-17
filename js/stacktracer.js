/** @jsx React.DOM */
var Output = React.createClass({
    render: function() {
        var outputRows = this.props.serverOutput.map(function(out) {
            return (
                <p><a href={out.link}>{out.filename} at line {out.linenum}</a></p>
            );
        });
        return (
            <div className="output">
            {outputRows}
            </div>
        );
    }
});
var Input = React.createClass({
    onTrace: function() {
        var trace = this.refs.userInput.getDOMNode().value.trim();
        this.props.handleTrace(trace);
    },
    render: function() {
        return (
            <div className="input" style={{float: "left"}}>
            <button onClick={this.onTrace}>trace my stack!</button><br />
            <textarea ref="userInput" rows="30" cols="80" style={{height: 'auto'}}></textarea>
            </div>
        )
    }
});
var Tracer = React.createClass({
    getInitialState: function() {
        return {serverOutput: []}
    },
    handleTrace: function(trace) {
        $.ajax({
            url: "parse",
            dataType: "json",
            type: "POST",
            data: {trace: trace},
            success: function(output) {
                if (output !== null) {
                    this.setState({serverOutput: output})
                } else {
                    this.setState({serverOutput: []})
                }
            }.bind(this),
            error: function(xhr, status, err) {
                console.error("parse", status, err.toString());
            }.bind(this),
        });
    },
    render: function() {
        return (
            <div className="tracer">
            <Input handleTrace={this.handleTrace}/>
            <Output serverOutput={this.state.serverOutput}/>
            </div>
        );
    }
});
React.renderComponent(
    <Tracer pollInterval={2000} />, document.getElementById("tracer")
);
