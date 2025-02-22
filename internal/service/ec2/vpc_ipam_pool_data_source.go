package ec2

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceVPCIpamPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCIpamPoolRead,

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// computed
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicly_advertisable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allocation_default_netmask_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allocation_max_netmask_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allocation_min_netmask_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allocation_resource_tags": tftags.TagsSchema(),
			"auto_import": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aws_service": {
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
			},
			"ipam_scope_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipam_scope_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_ipam_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceVPCIpamPoolRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	input := &ec2.DescribeIpamPoolsInput{}

	if v, ok := d.GetOk("ipam_pool_id"); ok {
		input.IpamPoolIds = aws.StringSlice([]string{v.(string)})

	}

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		input.Filters = BuildFiltersDataSource(filters.(*schema.Set))
	}

	output, err := conn.DescribeIpamPools(input)
	var pool *ec2.IpamPool

	if err != nil {
		return err
	}

	if output == nil || len(output.IpamPools) == 0 || output.IpamPools[0] == nil {
		return nil
	}
	pool = output.IpamPools[0]

	d.SetId(aws.StringValue(pool.IpamPoolId))

	if pool.PubliclyAdvertisable != nil {
		d.Set("publicly_advertisable", pool.PubliclyAdvertisable)
	}
	scopeId := strings.Split(*pool.IpamScopeArn, "/")[1]

	d.Set("allocation_resource_tags", KeyValueTags(ec2TagsFromIpamAllocationTags(pool.AllocationResourceTags)).Map())
	d.Set("auto_import", pool.AutoImport)
	d.Set("arn", pool.IpamPoolArn)
	d.Set("description", pool.Description)
	d.Set("ipam_scope_id", scopeId)
	d.Set("ipam_scope_type", pool.IpamScopeType)
	d.Set("locale", pool.Locale)
	d.Set("pool_depth", pool.PoolDepth)
	d.Set("aws_service", pool.AwsService)
	d.Set("source_ipam_pool_id", pool.SourceIpamPoolId)
	d.Set("state", pool.State)

	return nil
}
