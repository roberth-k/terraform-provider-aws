package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/outposts"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAwsOutpostsOutpost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsOutpostsOutpostRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsOutpostsOutpostRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).outpostsconn

	input := &outposts.ListOutpostsInput{}

	var outposts []*outposts.Outpost

	for {
		output, err := conn.ListOutposts(input)

		if err != nil {
			return fmt.Errorf("error listing Outposts Outposts: %w", err)
		}

		if output == nil {
			break
		}

		for _, outpost := range output.Outposts {
			if outpost == nil {
				continue
			}

			if v, ok := d.GetOk("id"); ok && v.(string) != aws.StringValue(outpost.OutpostId) {
				continue
			}

			if v, ok := d.GetOk("name"); ok && v.(string) != aws.StringValue(outpost.Name) {
				continue
			}

			outposts = append(outposts, outpost)
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	if len(outposts) == 0 {
		return fmt.Errorf("no Outposts Outpost found matching criteria; try different search")
	}

	if len(outposts) > 1 {
		return fmt.Errorf("multiple Outposts Outpost found matching criteria; try different search")
	}

	outpost := outposts[0]

	d.SetId(aws.StringValue(outpost.OutpostId))
	d.Set("arn", outpost.OutpostArn)
	d.Set("availability_zone", outpost.AvailabilityZone)
	d.Set("availability_zone_id", outpost.AvailabilityZoneId)
	d.Set("description", outpost.Description)
	d.Set("name", outpost.Name)
	d.Set("owner_id", outpost.OwnerId)
	d.Set("site_id", outpost.SiteId)

	return nil
}
