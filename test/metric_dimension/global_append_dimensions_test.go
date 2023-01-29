// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package metric_dimension

import (
	"github.com/aws/amazon-cloudwatch-agent-test/test/metric"
	"github.com/aws/amazon-cloudwatch-agent-test/test/metric/dimension"
	"github.com/aws/amazon-cloudwatch-agent-test/test/status"
	"github.com/aws/amazon-cloudwatch-agent-test/test/test_runner"
	"log"
)

type GlobalAppendDimensionsTestRunner struct {
	test_runner.BaseTestRunner
}

var _ test_runner.ITestRunner = (*GlobalAppendDimensionsTestRunner)(nil)

func (t *GlobalAppendDimensionsTestRunner) Validate() status.TestGroupResult {
	metricsToFetch := t.GetMeasuredMetrics()
	testResults := make([]status.TestResult, len(metricsToFetch))
	for i, metricName := range metricsToFetch {
		testResults[i] = t.validateNoAppendDimensionMetric(metricName)
	}

	return status.TestGroupResult{
		Name:        t.GetTestName(),
		TestResults: testResults,
	}
}

func (t *GlobalAppendDimensionsTestRunner) GetTestName() string {
	return "GlobalAppendDimension"
}

func (t *GlobalAppendDimensionsTestRunner) GetAgentConfigFileName() string {
	return "global_append_dimension.json"
}

func (t *GlobalAppendDimensionsTestRunner) GetMeasuredMetrics() []string {
	return []string{"cpu_time_active"}
}

func (t *GlobalAppendDimensionsTestRunner) validateNoAppendDimensionMetric(metricName string) status.TestResult {
	testResult := status.TestResult{
		Name:   metricName,
		Status: status.FAILED,
	}

	expDims, failed := t.DimensionFactory.GetDimensions([]dimension.Instruction{
		{
			Key:   "ImageId",
			Value: dimension.UnknownDimensionValue(),
		},
		{
			Key:   "InstanceId",
			Value: dimension.UnknownDimensionValue(),
		},
		{
			Key:   "InstanceType",
			Value: dimension.UnknownDimensionValue(),
		},
	})

	if len(failed) > 0 {
		return testResult
	}

	fetcher := metric.MetricValueFetcher{}
	values, err := fetcher.Fetch("MetricAppendDimensionTest", metricName, expDims, metric.AVERAGE)
	log.Printf("metric values are %v", values)
	if err != nil {
		return testResult
	}

	if !isAllValuesGreaterThanOrEqualToZero(metricName, values) {
		return testResult
	}

	dropDims, failed := t.DimensionFactory.GetDimensions([]dimension.Instruction{
		{
			Key:   "host",
			Value: dimension.UnknownDimensionValue(),
		},
	})

	if len(failed) > 0 {
		return testResult
	}

	values, err = fetcher.Fetch("MetricAppendDimensionTest", metricName, dropDims, metric.AVERAGE)
	if err != nil || len(values) != 0 {
		return testResult
	}

	testResult.Status = status.SUCCESSFUL
	return testResult
}
